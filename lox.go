package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/sanity-io/litter"
)

// Lox encapsulates program state and execution methods
type Lox struct {
	HadError bool
}

func (l *Lox) runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	source := string(bytes)
	l.run(source)

	if (l.HadError) {
		os.Exit(65)
	}
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		l.run(text)
	}
}

func (l *Lox) run(source string) {
	scanner := makeScanner(source)
	tokens := scanner.scanTokens()
	litter.Dump(tokens)
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	msg := fmt.Sprintf("[line %d ] Error%s : %s", line, where, message)
	fmt.Fprintln(os.Stderr, msg)
	l.HadError = true
}
