package main

import (
	"os"
	"fmt"
)

var lox *Lox

func main() {
	lox = &Lox{}

	args := os.Args[1:]
	if (len(args) > 1) {                                   
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if (len(args) == 1) {
		lox.runFile(args[0])
	} else {                                                 
		lox.runPrompt()
	} 
}