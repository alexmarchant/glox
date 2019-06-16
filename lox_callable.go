package main

type LoxCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, *RuntimeError)
	Arity() int
}
