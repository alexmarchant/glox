package main

import "fmt"

type LoxFunction struct {
	Declaration   *FunctionStmt
	Closure       *Environment
	IsInitializer bool
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
		// Short circuit return statements using errors :P
		if err.Return != nil {
			// Return instance instead of nil from `return;` in
			// init methods
			if f.IsInitializer {
				f.Closure.getAt(0, "this")
			}

			return err.Return, nil
		} else {
			return nil, err
		}
	}

	// Return instance of class implicetly from init methods
	if f.IsInitializer {
		return f.Closure.getAt(0, "this")
	}

	return nil, nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

func (f *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.Closure)
	environment.define("this", instance)
	return &LoxFunction{
		Declaration:   f.Declaration,
		Closure:       environment,
		IsInitializer: f.IsInitializer,
	}
}
