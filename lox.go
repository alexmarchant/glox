package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Lox struct {
	Interpreter     *Interpreter
	HadError        bool
	HadRuntimeError bool
}

func NewLox() *Lox {

	return &Lox{
		Interpreter: NewInterpreter(),
	}
}

func (l *Lox) runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	source := string(bytes)
	l.run(source)

	if l.HadError {
		os.Exit(65)
	}

	if l.HadRuntimeError {
		os.Exit(70)
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
		l.resetErrorState()
	}
}

func (l *Lox) run(source string) {
	scanner := makeScanner(source)
	tokens := scanner.scanTokens()
	parser := &Parser{Tokens: tokens}
	statements := parser.parse()

	if lox.HadError {
		return
	}

	l.Interpreter.Interpret(statements)
}

func (l *Lox) errorLine(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) errorToken(token *Token, message string) {
	if token.Type == EOF {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (l *Lox) report(line int, where string, message string) {
	msg := fmt.Sprintf("[line %d] Error%s : %s", line, where, message)
	fmt.Fprintln(os.Stderr, msg)
	l.HadError = true
}

func (l *Lox) runtimeError(err *RuntimeError) {
	msg := fmt.Sprintf("%s\n[line %d]", err.Message, err.Token.Line)
	fmt.Fprintln(os.Stderr, msg)
	l.HadRuntimeError = true
}

func (l *Lox) resetErrorState() {
	l.HadError = false
	l.HadRuntimeError = false
}
