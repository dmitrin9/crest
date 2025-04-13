package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getHelpString() string {
	helpString := ""
	helpString += "run                Run a Crestfile.\n"
	helpString += "-t/--test-http     Test http mode.\n"
	helpString += "-v/--verbose       Print in verbose mode.\n"
	helpString += "-q/--quiet         Print in quiet mode.\n"
	helpString += "-f/--follow-robots Follow robots.txt.\n"
	return helpString
}

func getPathToCrestfile(args []string) string {
	pathToCrestfile := strings.Split(args[len(args)-1], "/")
	crestfileName := strings.ToLower(pathToCrestfile[len(pathToCrestfile)-1])
	return crestfileName
}

func main() {
	ctx := Context{
		verbose:      false,
		quiet:        false,
		followRobots: false,

		exclude: []string{},
	}
	s := State{
		raw:         "",
		lexNodes:    []LexNode{},
		parserNodes: []ParserNode{},

		variable:       make(map[string]string),
		instructionSet: []string{},

		offset: 0,
		row:    0,
		col:    0,
	}
	args := os.Args
	if len(args) >= 2 {
		if args[1] == "run" && getPathToCrestfile(args) == "crestfile" {
			if err := HandleFile(args, &s, &ctx); err != nil {
				log.Fatal(err)
			}
		} else if args[1] == "run" {
			fmt.Fprintln(os.Stderr, "It seems like you inputted an invalid path for your Crestfile.")
		} else {
			if err := Handle(args, &ctx); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, getHelpString())
	}
}
