package main

import (
	"fmt"
)

type Interpreter struct {
	Error       *RuntimeError
	Globals     *Environment
	Environment *Environment
	Locals      map[Expr]int
}

type RuntimeError struct {
	Token   *Token
	Message string
	Return  interface{}
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment(nil)

	// Native functions
	env.define("clock", &ClockNativeFunc{})
	env.define("print", &PrintNativeFunc{})

	return &Interpreter{
		Environment: env,
		Globals:     env,
		Locals:      map[Expr]int{},
	}
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			lox.runtimeError(err)
			return
		}
	}
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.Locals[expr] = depth
}

func (i *Interpreter) lookupVariable(name *Token, expr Expr) (interface{}, *RuntimeError) {
	if distance, ok := i.Locals[expr]; ok {
		return i.Environment.getAt(distance, name.Lexeme)
	} else {
		return i.Globals.get(name)
	}
}

// Statements

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) (interface{}, *RuntimeError) {
	_, err := i.evaluate(stmt.Expression)
	return nil, err
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) (interface{}, *RuntimeError) {
	var value interface{}
	var err *RuntimeError

	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.Environment.define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) (interface{}, *RuntimeError) {
	val, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	for i.isTruthy(val) {
		err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
		val, err = i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) (interface{}, *RuntimeError) {
	i.executeBlock(
		stmt.Statements,
		NewEnvironment(i.Environment))
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) (interface{}, *RuntimeError) {
	val, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(val) {
		err = i.execute(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		err = i.execute(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) (interface{}, *RuntimeError) {
	function := &LoxFunction{
		Declaration: stmt,
		Closure:     i.Environment,
	}
	i.Environment.define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) (interface{}, *RuntimeError) {
	var value interface{}
	var err *RuntimeError

	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, &RuntimeError{Return: value}
}

func (i *Interpreter) VisitClassStmt(stmt *ClassStmt) (interface{}, *RuntimeError) {
	var superclass *LoxClass

	if stmt.Superclass != nil {
		var ok bool

		val, err := i.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}
		superclass, ok = val.(*LoxClass)
		if !ok {
			return nil, &RuntimeError{
				Token:   stmt.Superclass.Name,
				Message: "Superclass must be a class.",
			}
		}
	}

	i.Environment.define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.Environment = NewEnvironment(i.Environment)
		i.Environment.define("super", superclass)
	}

	methods := map[string]*LoxFunction{}
	for _, method := range stmt.Methods {
		function := &LoxFunction{
			Declaration:   method,
			Closure:       i.Environment,
			IsInitializer: method.Name.Lexeme == "init",
		}
		methods[method.Name.Lexeme] = function
	}

	class := &LoxClass{
		Name:       stmt.Name.Lexeme,
		Superclass: superclass,
		Methods:    methods,
	}

	if superclass != nil {
		i.Environment = i.Environment.Enclosing
	}

	i.Environment.assign(stmt.Name, class)
	return nil, nil
}

// Expressions

func (i *Interpreter) VisitGetExpr(expr *GetExpr) (interface{}, *RuntimeError) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if instance, ok := object.(*LoxInstance); ok {
		return instance.get(expr.Name)
	}

	return nil, &RuntimeError{
		Token:   expr.Name,
		Message: "Only instances have properties.",
	}
}

func (i *Interpreter) VisitSetExpr(expr *SetExpr) (interface{}, *RuntimeError) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*LoxInstance)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Name,
			Message: "Only instances have properties.",
		}
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	instance.set(expr.Name, value)
	return value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, *RuntimeError) {
	val, err := i.evaluate(expr.Expression)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, *RuntimeError) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, *RuntimeError) {
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

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, *RuntimeError) {
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
		return nil, &RuntimeError{
			Token:   expr.Operator,
			Message: msg,
		}
	case Slash:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if right.(float64) == 0 {
			msg := "Cannot divide by 0."
			return nil, &RuntimeError{
				Token:   expr.Operator,
				Message: msg,
			}
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

func (i *Interpreter) VisitVarExpr(expr *VarExpr) (interface{}, *RuntimeError) {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) (interface{}, *RuntimeError) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if distance, ok := i.Locals[expr]; ok {
		i.Environment.assignAt(distance, expr.Name, value)
	} else {
		err = i.Globals.assign(expr.Name, value)
	}
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) (interface{}, *RuntimeError) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == Or {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) (interface{}, *RuntimeError) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []interface{}{}
	for _, arg := range expr.Arguments {
		res, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, res)
	}

	// Cast as callable
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Paren,
			Message: "Can only call functions and classes.",
		}
	}

	// Check arity
	if len(arguments) != function.Arity() {
		msg := fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))
		return nil, &RuntimeError{
			Token:   expr.Paren,
			Message: msg,
		}
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitThisExpr(expr *ThisExpr) (interface{}, *RuntimeError) {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) (interface{}, *RuntimeError) {
	distance := i.Locals[expr]
	superclassInterface, err := i.Environment.getAt(distance, "super")
	if err != nil {
		return nil, err
	}
	superclass := superclassInterface.(*LoxClass)

	instanceInterface, err := i.Environment.getAt(distance-1, "this")
	if err != nil {
		return nil, err
	}
	instance := instanceInterface.(*LoxInstance)

	method, ok := superclass.findMethod(expr.Method.Lexeme)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Method,
			Message: fmt.Sprintf("Undefined property '%s'.", expr.Method.Lexeme),
		}
	}
	return method.bind(instance), nil
}

// Helpers

func (i *Interpreter) execute(stmt Stmt) *RuntimeError {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) *RuntimeError {
	previousEnv := i.Environment
	i.Environment = env
	defer func() {
		i.Environment = previousEnv
	}()

	for _, stmt := range statements {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, *RuntimeError) {
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

func (i *Interpreter) checkNumberOperand(operator *Token, operand interface{}) *RuntimeError {
	_, ok := operand.(float64)
	if ok {
		return nil
	}
	return &RuntimeError{
		Token:   operator,
		Message: "Operand must be number.",
	}
}

func (i *Interpreter) checkNumberOperands(operator *Token, left interface{}, right interface{}) *RuntimeError {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return nil
	}
	return &RuntimeError{
		Token:   operator,
		Message: "Operands must be numbers.",
	}
}

func (i *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}

	if fVal, ok := val.(float64); ok {
		sVal := fmt.Sprintf("%f", fVal)
		return sVal
	}

	return fmt.Sprintf("%v", val)
}
