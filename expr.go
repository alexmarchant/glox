package main

type ExprType int
const (
	ExprTypeUnaryExpr ExprType = iota
	ExprTypeBinaryExpr
	ExprTypeGroupingExpr
	ExprTypeLiteralExpr
)

type Expr interface {
	ExprType() ExprType
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) (interface{}, error)
	VisitGroupingExpr(*GroupingExpr) (interface{}, error)
	VisitLiteralExpr(*LiteralExpr) (interface{}, error)
	VisitUnaryExpr(*UnaryExpr) (interface{}, error)
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
	Value LiteralValue
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

