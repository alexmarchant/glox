package main

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		Values: map[string]interface{}{},
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) get(name *Token) (interface{}, *RuntimeError) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}
	return nil, &RuntimeError{
		Token:   name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}
