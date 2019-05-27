package main

type ExprType int
const (
	ExprTypeGroupingExpr ExprType = iota
	ExprTypeLiteralExpr
	ExprTypeUnaryExpr
	ExprTypeBinaryExpr
)

type Expr interface {
	ExprType() ExprType
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitUnaryExpr(*UnaryExpr) (interface{}, error)
	VisitBinaryExpr(*BinaryExpr) (interface{}, error)
	VisitGroupingExpr(*GroupingExpr) (interface{}, error)
	VisitLiteralExpr(*LiteralExpr) (interface{}, error)
}

type GroupingExpr struct {
	Expression Expr
}

func (t *GroupingExpr) ExprType() ExprType {
	return ExprTypeGroupingExpr
}

func (t *GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(t)
}

type LiteralExpr struct {
	Value interface{}
}

func (t *LiteralExpr) ExprType() ExprType {
	return ExprTypeLiteralExpr
}

func (t *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(t)
}

type UnaryExpr struct {
	Operator *Token
	Right Expr
}

func (t *UnaryExpr) ExprType() ExprType {
	return ExprTypeUnaryExpr
}

func (t *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(t)
}

type BinaryExpr struct {
	Left Expr
	Operator *Token
	Right Expr
}

func (t *BinaryExpr) ExprType() ExprType {
	return ExprTypeBinaryExpr
}

func (t *BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(t)
}

