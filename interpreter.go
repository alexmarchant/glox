package main

import (
	"fmt"
	"errors"
	"strings"
)

type Interpreter struct {
	Error *RuntimeError
}

type RuntimeError struct {
	Token *Token
	Message string
}

func (i *Interpreter) Interpret(expr Expr) {
	res, err := i.evaluate(expr)
	if err != nil {
		lox.runtimeError(i.Error)
		return
	}
	fmt.Printf("=> %s\n", i.stringify(res))
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	val, err := i.evaluate(expr.Expression)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Bang:
		return i.isTruthy(right), nil
	case Minus:
		return -(right.(float64)), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case BangEqual:
		return !i.isEqual(left, right), nil
	case EqualEqual:
		return i.isEqual(left, right), nil
	case Greater:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case Less:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case Minus:
		err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case Plus:
		// Number
		fltLeft, isFltLeft := left.(float64)
		fltRight, isFltRight := right.(float64)
		if isFltLeft && isFltRight {
			return fltLeft + fltRight, nil
		}

		// String
		strLeft, isStrLeft := left.(string)
		strRight, isStrRight := right.(string)
		if isStrLeft && isStrRight {
			return strLeft + strRight, nil
		}

		msg := "Operands must be two numbers or two strings."
		i.Error = &RuntimeError{
			Token: expr.Operator,
			Message: msg,
		}
		return nil, errors.New(msg)
	case Slash:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if right.(float64) == 0 {
			msg := "Cannot divide by 0."
			i.Error = &RuntimeError{
				Token: expr.Operator,
				Message: msg,
			}
			return nil, errors.New(msg)
		}
		return left.(float64) / right.(float64), nil
	case Star:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	return nil, nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if bVal, ok := val.(bool); ok {
		return bVal
	}
	return true
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) checkNumberOperand(operator *Token, operand interface{}) error {
	_, ok := operand.(float64)
	if ok {
		return nil
	}
	msg := "Operand must be number."
	i.Error = &RuntimeError{
		Token: operator,
		Message: msg,
	}
	return errors.New(msg)
}

func (i *Interpreter) checkNumberOperands(operator *Token, left interface{}, right interface{}) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return nil
	}
	msg := "Operands must be numbers."
	i.Error = &RuntimeError{
		Token: operator,
		Message: msg,
	}
	return errors.New(msg)
}

func (i *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}

	if fVal, ok := val.(float64); ok {
		sVal := fmt.Sprintf("%f", fVal)
		sVal = strings.Trim(sVal, "0")
		sVal = strings.Trim(sVal, ".")
		return sVal
	}

	return fmt.Sprintf("%v", val)
}
