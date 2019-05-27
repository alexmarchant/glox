package main

import (
	"fmt"
)

type Token struct {
	Type TokenType
	Lexeme string
	Literal LiteralValue
	Line int
}

func (t *Token) String() string {
	return fmt.Sprintf("<Token type: %s, lexeme: %s, literal: %v>", t.Type, t.Lexeme, t.Literal)
}

type LiteralValue interface {
	LiteralValueType() LiteralValueType
	String() string
}

type LiteralValueType int
const (
	LiteralValueTypeString LiteralValueType = iota
	LiteralValueTypeNumber
	LiteralValueTypeBool
	LiteralValueTypeNil
)

type LiteralValueString struct {
	Value string
}

func (l *LiteralValueString) LiteralValueType() LiteralValueType {
	return LiteralValueTypeString
}

func (l *LiteralValueString) String() string {
	return l.Value
}

type LiteralValueNumber struct {
	Value float64
}

func (l *LiteralValueNumber) LiteralValueType() LiteralValueType {
	return LiteralValueTypeNumber
}

func (l *LiteralValueNumber) String() string {
	return fmt.Sprintf("%f", l.Value)
}

type LiteralValueBool struct {
	Value bool
}

func (l *LiteralValueBool) LiteralValueType() LiteralValueType {
	return LiteralValueTypeBool
}

func (l *LiteralValueBool) String() string {
	return fmt.Sprintf("%t", l.Value)
}

type LiteralValueNil struct {}

// LiteralValueType is the type of literal value
func (l *LiteralValueNil) LiteralValueType() LiteralValueType {
	return LiteralValueTypeNil
}

func (l *LiteralValueNil) String() string {
	return "nil"
}