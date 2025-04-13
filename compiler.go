package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var TOKS map[string]string = map[string]string{
	"=":            "ASSIGNMENT",
	"true":         "BOOL",
	"false":        "BOOL",
	"exclude":      "SET",
	"verbose":      "SET",
	"quiet":        "SET",
	"followRobots": "SET",
	"type":         "SET",
	"url":          "SET",
	"testHTTP":     "TEST_TYPE",
}

/*
 * Nodes {
 * 	LexNode {
 *    tok_raw: "url"
 *    tok_type: "KEYWORD"
 *  },
 *  LexNode {
 *    tok_raw: "=",
 *    tok_type: "ASSIGNMENT"
 *  },
 *  LexNode {
 *    tok_raw: "example.com",
 *    tok_type: "STRING"
 *  }
 * }
 */
type LexNode struct {
	tok_raw  string
	tok_type string
}

func newLexnode(raw string, tokType string) LexNode {
	return LexNode{
		tok_raw:  raw,
		tok_type: tokType,
	}
}

/*
 * ParserNode {
 * 	operation: Assignment
 *  operands:{
 *	  LexNode {
 *		 tok_raw: url
 * 		 tok_type: nil
 * 	 },
 *	 LexNode {
 * 		tok_raw: "example.com"
 *		tok_type: STRING
 *	 }
 * }
 */
type ParserNode struct {
	operation string
	operands  []LexNode
}

type State struct {
	raw         string
	lexNodes    []LexNode
	parserNodes []ParserNode

	variable map[string]string

	instructionSet []string

	offset int
	row    int
	col    int
}

func (a *ParserNode) Clear() {
	a.operation = ""
	a.operands = nil
}

func (s *State) lexError(message string) error {
	errorMessage := fmt.Sprintf("LEXER ERROR: %s at row:%d, col:%d", message, s.row, s.col)
	return errors.New(errorMessage)
}

func (s *State) parseError(message string, operation string) error {
	errorMessage := fmt.Sprintf("PARSER ERROR: %s at operation %s", message, operation)
	return errors.New(errorMessage)
}

func (s *State) compileError(message string, operation string) error {
	errorMessage := fmt.Sprintf("COMPILER ERROR: %s during operation %s", message, operation)
	return errors.New(errorMessage)
}

func (s *State) next() {
	if s.offset == len(s.raw) {
		return
	} else if s.raw[s.offset] == '\n' {
		s.col++
		s.row = 0
		s.offset++
	} else {
		s.row++
		s.offset++
	}
}

func (s *State) Lexer() error {
	var nodes []LexNode

	if len(s.raw) == 0 {
		return s.lexError("Empty file data")
	}

	re := regexp.MustCompile("\"[^\"]*\"|`[^`]*`|\\S+")
	splitted := re.FindAllString(s.raw, -1)

	for _, _ = range splitted {
		i := s.offset
		c := splitted[i]

		var token LexNode
		rawToken := c

		if rawToken[0] == '"' && rawToken[len(rawToken)-1] == '"' {
			token = newLexnode(rawToken, "STRING")
		} else if rawToken[0] == '`' && rawToken[len(rawToken)-1] == '`' {
			token = newLexnode(rawToken, "CODE")
		} else if rawToken[0] == '{' && rawToken[len(rawToken)-1] == '}' {
			token = newLexnode(rawToken, "VARIABLE")
		} else if len(TOKS[rawToken]) > 0 {
			token = newLexnode(rawToken, TOKS[rawToken])
		} else {
			token = newLexnode(rawToken, "")
		}

		nodes = append(nodes, token)
		s.next()
	}

	s.lexNodes = nodes
	return nil
}

func (s *State) Parser() error {
	tokens := s.lexNodes
	var parserNodes []ParserNode

	for i, c := range tokens {
		var node ParserNode

		if c.tok_type == "ASSIGNMENT" {
			node.operation = "Assignment"
			if i > 0 && i < len(tokens)-1 {
				name := tokens[i-1]
				value := tokens[i+1]
				node.operands = []LexNode{name, value}
			} else {
				return s.parseError("Assignment failed because assignment operator is in an invalid location", node.operation)
			}
			parserNodes = append(parserNodes, node)
			node.Clear()
		}
		if c.tok_type == "SET" {
			node.operation = "Set"
			if i < len(tokens)-1 {
				name := tokens[i]
				value := tokens[i+1]
				node.operands = []LexNode{name, value}
			} else {
				return s.parseError("Set operator is in an invalid location", node.operation)
			}
			parserNodes = append(parserNodes, node)
			node.Clear()
		}
	}

	s.parserNodes = parserNodes

	return nil
}

func (s *State) compileCode(codeblock string) string {
	var varBuf string
	var valueBuf string
	var insideVariable bool
	var changeIdx int
	re := regexp.MustCompile("\\s")
	split := re.Split(strings.TrimSpace(codeblock), -1)

	for i, c := range split {
		varBuf = ""
		valueBuf = ""
		insideVariable = false
		changeIdx = -1
		for _, d := range c {
			if d == '{' {
				insideVariable = true
				continue
			} else if d == '}' {
				insideVariable = false
			}

			if insideVariable == true {
				varBuf += string(d)
			} else if len(varBuf) > 0 {
				variable := varBuf
				valueBuf += s.variable[variable]
				changeIdx = i
				varBuf = ""
				fmt.Println("VALUE BUFFER ", valueBuf)
			}
		}
		if changeIdx != -1 {
			split[i] = valueBuf
			changeIdx = -1
		}
	}
	code := strings.Join(split, " ")
	return code
}

func (s *State) Compiler() error {
	/*
	 * Compiles the parse nodes to a simple instruction set
	 * 		( this also invovles substituting variable names
	 *        and simplifying code blocks )
	 * which can be further interpreted by a function
	 * in the crest.go file to generate a runtime.
	 */
	s.variable = make(map[string]string)

	for _, c := range s.parserNodes {
		if c.operation == "Assignment" {
			var name string
			var value string

			nameToken := c.operands[0]
			valueToken := c.operands[1]
			if nameToken.tok_type != "" {
				s.compileError("Invalid type for variable name", c.operation)
			}
			name = nameToken.tok_raw
			value = valueToken.tok_raw

			if valueToken.tok_type == "STRING" {
				value = value[1 : len(value)-1]
			}
			if valueToken.tok_raw == "CURRENT" {
				value = "CURRENT"
			}
			if valueToken.tok_raw == "CONTENT" {
				value = "CONTENT"
			}
			s.variable[name] = value
		}
		if c.operation == "Set" {
			var name string
			var value string

			nameToken := c.operands[0]
			valueToken := c.operands[1]
			if nameToken.tok_type != "" {
				s.compileError("Invalid type for set name", c.operation)
			}
			if valueToken.tok_type == "VARIABLE" {
				variableValue := ""
				rawVariable := valueToken.tok_raw[1 : len(valueToken.tok_raw)-1]
				variableValue = s.variable[rawVariable]
				if len(variableValue) <= 0 {
					return s.compileError("Variable not found", c.operation)
				}
				value = variableValue
			} else if valueToken.tok_type == "CODE" {
				value = s.compileCode(valueToken.tok_raw[1 : len(valueToken.tok_raw)-1])
			} else if valueToken.tok_type == "STRING" {
				value = valueToken.tok_raw[1 : len(valueToken.tok_raw)-1]
			} else {
				value = valueToken.tok_raw
			}
			name = nameToken.tok_raw
			s.instructionSet = append(s.instructionSet, name)
			s.instructionSet = append(s.instructionSet, value)
		}
	}
	return nil
}
