package main

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    map[string]interface{}{},
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) assign(name *Token, value interface{}) *RuntimeError {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}

	return &RuntimeError{
		Token:   name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

func (e *Environment) get(name *Token) (interface{}, *RuntimeError) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.get(name)
	}

	return nil, &RuntimeError{
		Token:   name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}
