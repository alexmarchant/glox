package main

type Stmt interface {
	Accept(StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) (interface{}, error)
	VisitPrintStmt(*PrintStmt) (interface{}, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (t *ExpressionStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(t)
}

type PrintStmt struct {
	Expression Expr
}

func (t *PrintStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(t)
}

