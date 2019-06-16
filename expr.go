package main

type Expr interface {
	Accept(ExprVisitor) (interface{}, *RuntimeError)
}

type ExprVisitor interface {
	VisitLiteralExpr(*LiteralExpr) (interface{}, *RuntimeError)
	VisitVarExpr(*VarExpr) (interface{}, *RuntimeError)
	VisitGetExpr(*GetExpr) (interface{}, *RuntimeError)
	VisitThisExpr(*ThisExpr) (interface{}, *RuntimeError)
	VisitBinaryExpr(*BinaryExpr) (interface{}, *RuntimeError)
	VisitGroupingExpr(*GroupingExpr) (interface{}, *RuntimeError)
	VisitLogicalExpr(*LogicalExpr) (interface{}, *RuntimeError)
	VisitCallExpr(*CallExpr) (interface{}, *RuntimeError)
	VisitSetExpr(*SetExpr) (interface{}, *RuntimeError)
	VisitUnaryExpr(*UnaryExpr) (interface{}, *RuntimeError)
	VisitAssignExpr(*AssignExpr) (interface{}, *RuntimeError)
}

type LogicalExpr struct {
	Left Expr
	Operator *Token
	Right Expr
}

func (t *LogicalExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitLogicalExpr(t)
}

type CallExpr struct {
	Callee Expr
	Paren *Token
	Arguments []Expr
}

func (t *CallExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitCallExpr(t)
}

type SetExpr struct {
	Object Expr
	Name *Token
	Value Expr
}

func (t *SetExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitSetExpr(t)
}

type UnaryExpr struct {
	Operator *Token
	Right Expr
}

func (t *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitUnaryExpr(t)
}

type AssignExpr struct {
	Name *Token
	Value Expr
}

func (t *AssignExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitAssignExpr(t)
}

type LiteralExpr struct {
	Value interface{}
}

func (t *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitLiteralExpr(t)
}

type VarExpr struct {
	Name *Token
}

func (t *VarExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitVarExpr(t)
}

type GetExpr struct {
	Object Expr
	Name *Token
}

func (t *GetExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitGetExpr(t)
}

type ThisExpr struct {
	Keyword *Token
}

func (t *ThisExpr) Accept(visitor ExprVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitThisExpr(t)
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

