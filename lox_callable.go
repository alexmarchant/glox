package main

import (
	"fmt"
	"time"
)

type LoxCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, *RuntimeError)
	Arity() int
}

type LoxFunction struct {
	Declaration *FunctionStmt
	Closure     *Environment
}

func (f *LoxFunction) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	// Setup scope
	environment := NewEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		environment.define(param.Lexeme, args[i])
	}

	// Execute body
	err := i.executeBlock(f.Declaration.Body, environment)
	if err != nil {
		if err.Return != nil {
			return err.Return, nil
		} else {
			return nil, err
		}
	}

	return nil, nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

// print()
type PrintNativeFunc struct{}

func (f *PrintNativeFunc) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	fmt.Printf("%s\n", i.stringify(args[0]))
	return nil, nil
}

func (f *PrintNativeFunc) Arity() int {
	return 1
}

func (f *PrintNativeFunc) String() string {
	return "<native fn>"
}

// clock()
type ClockNativeFunc struct{}

func (f *ClockNativeFunc) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	return time.Now().UTC().UnixNano() / 1000000000, nil
}

func (f *ClockNativeFunc) Arity() int {
	return 0
}

func (f *ClockNativeFunc) String() string {
	return "<native fn>"
}
