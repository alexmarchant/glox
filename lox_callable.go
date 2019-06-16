package main

import (
	"fmt"
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
