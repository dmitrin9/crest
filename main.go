package main

import (
	"log"
	"os"
)

func main() {
	ctx := Context{
		verbose:      false,
		quiet:        false,
		followRobots: false,
		hooks:        make(map[string]string),

		exclude: []string{},
	}
	s := State{
		raw:         "",
		lexNodes:    []LexNode{},
		parserNodes: []ParserNode{},

		variable:       make(map[string]string),
		hooks:          make(map[string]string),
		instructionSet: []string{},

		offset: 0,
		row:    0,
		col:    0,
	}
	args := os.Args
	if len(args) >= 2 {
		if args[1] == "run" {
			if err := HandleFile(args, &s, &ctx); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := Handle(args, &ctx); err != nil {
				log.Fatal(err)
			}
		}
	}
}
