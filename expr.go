package main

type Expr interface {
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitLiteralExpr(*LiteralExpr) (interface{}, error)
	VisitUnaryExpr(*UnaryExpr) (interface{}, error)
	VisitBinaryExpr(*BinaryExpr) (interface{}, error)
	VisitGroupingExpr(*GroupingExpr) (interface{}, error)
}

type BinaryExpr struct {
	Left Expr
	Operator *Token
	Right Expr
}

func (t *BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(t)
}

type GroupingExpr struct {
	Expression Expr
}

func (t *GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(t)
}

type LiteralExpr struct {
	Value interface{}
}

func (t *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(t)
}

type UnaryExpr struct {
	Operator *Token
	Right Expr
}

func (t *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(t)
}
