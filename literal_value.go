package main

type LiteralValue interface {
	LiteralValueType() LiteralValueType
	Value() interface{}
}

type LiteralValueType int
const (
	LiteralValueTypeString LiteralValueType = iota
	LiteralValueTypeNumber
	LiteralValueTypeBool
	LiteralValueTypeNil
)

type LiteralValueString struct {
	StringValue string
}

func (l *LiteralValueString) LiteralValueType() LiteralValueType {
	return LiteralValueTypeString
}

func (l *LiteralValueString) Value() interface{} {
	return l.StringValue
}

type LiteralValueNumber struct {
	NumberValue float64
}

func (l *LiteralValueNumber) LiteralValueType() LiteralValueType {
	return LiteralValueTypeNumber
}

func (l *LiteralValueNumber) Value() interface{} {
	return l.NumberValue
}

type LiteralValueBool struct {
	BoolValue bool
}

func (l *LiteralValueBool) LiteralValueType() LiteralValueType {
	return LiteralValueTypeBool
}

func (l *LiteralValueBool) Value() interface{} {
	return l.BoolValue
}

type LiteralValueNil struct {}

func (l *LiteralValueNil) LiteralValueType() LiteralValueType {
	return LiteralValueTypeNil
}

func (l *LiteralValueNil) Value() interface{} {
	return nil
}
