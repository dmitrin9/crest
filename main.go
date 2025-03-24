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
	}
	args := os.Args
	if err := Handle(args, &ctx); err != nil {
		log.Fatal(err)
	}
}
