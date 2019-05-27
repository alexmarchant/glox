package main

type StmtType int
const (
	StmtTypePrintStmt StmtType = iota
	StmtTypeExpressionStmt
)

type Stmt interface {
	StmtType() StmtType
	Accept(StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) (interface{}, error)
	VisitPrintStmt(*PrintStmt) (interface{}, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (t *ExpressionStmt) StmtType() StmtType {
	return StmtTypeExpressionStmt
}

func (t *ExpressionStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(t)
}

type PrintStmt struct {
	Expression Expr
}

func (t *PrintStmt) StmtType() StmtType {
	return StmtTypePrintStmt
}

func (t *PrintStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(t)
}

