package main

import "fmt"

type LoxInstance struct {
	Class  *LoxClass
	Fields map[string]interface{}
}

func (l *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.Class.Name)
}

func (l *LoxInstance) get(name *Token) (interface{}, *RuntimeError) {
	if val, ok := l.Fields[name.Lexeme]; ok {
		return val, nil
	}

	if method, ok := l.Class.findMethod(name.Lexeme); ok {
		return method.bind(l), nil
	}

	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
	}
}

func (l *LoxInstance) set(name *Token, value interface{}) {
	l.Fields[name.Lexeme] = value
}
