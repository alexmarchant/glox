package main

type Expr interface {
	Accept(ExprVisitor) (interface{}, *RuntimeError)
}

type ExprVisitor interface {
	VisitLogicalExpr(*LogicalExpr) (interface{}, *RuntimeError)
	VisitCallExpr(*CallExpr) (interface{}, *RuntimeError)
	VisitBinaryExpr(*BinaryExpr) (interface{}, *RuntimeError)
	VisitGroupingExpr(*GroupingExpr) (interface{}, *RuntimeError)
	VisitLiteralExpr(*LiteralExpr) (interface{}, *RuntimeError)
	VisitUnaryExpr(*UnaryExpr) (interface{}, *RuntimeError)
	VisitVarExpr(*VarExpr) (interface{}, *RuntimeError)
	VisitAssignExpr(*AssignExpr) (interface{}, *RuntimeError)
}

type CallExpr struct {
	Callee Expr
	Paren *Token
	Arguments []Expr
}

func (t *CallExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitCallExpr(t)
}

type BinaryExpr struct {
	Left Expr
	Operator *Token
	Right Expr
}

func (t *BinaryExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitBinaryExpr(t)
}

type GroupingExpr struct {
	Expression Expr
}

func (t *GroupingExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitGroupingExpr(t)
}

type LiteralExpr struct {
	Value interface{}
}

func (t *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitLiteralExpr(t)
}

type UnaryExpr struct {
	Operator *Token
	Right Expr
}

func (t *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitUnaryExpr(t)
}

type VarExpr struct {
	Name *Token
}

func (t *VarExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitVarExpr(t)
}

type AssignExpr struct {
	Name *Token
	Value Expr
}

func (t *AssignExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitAssignExpr(t)
}

type LogicalExpr struct {
	Left Expr
	Operator *Token
	Right Expr
}

func (t *LogicalExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitLogicalExpr(t)
}

