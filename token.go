package main

import (
	"fmt"
)

// Token is a small key languange part
type Token struct {
	Type TokenType
	Lexeme string
	Literal LiteralValue
	Line int
}

func (t *Token) String() string {
	return fmt.Sprintf("<Token type: %s, lexeme: %s, literal: %v>", t.Type, t.Lexeme, t.Literal)
}

// LiteralValue is any struct which can hold the value of a token literal
type LiteralValue interface {
	LiteralValueType() LiteralValueType
	String() string
}

// LiteralValueType is an enum
type LiteralValueType int

const (
	LiteralValueTypeString LiteralValueType = iota
	LiteralValueTypeNumber
)

// LiteralValueString holds string literal val
type LiteralValueString struct {
	Value string
}

// LiteralValueType is the type of literal value
func (l *LiteralValueString) LiteralValueType() LiteralValueType {
	return LiteralValueTypeString
}

func (l *LiteralValueString) String() string {
	return l.Value
}

// LiteralValueNumber holds number literal val
type LiteralValueNumber struct {
	Value float64
}

// LiteralValueType is the type of literal value
func (l *LiteralValueNumber) LiteralValueType() LiteralValueType {
	return LiteralValueTypeNumber
}

func (l *LiteralValueNumber) String() string {
	return fmt.Sprintf("%f", l.Value)
}
