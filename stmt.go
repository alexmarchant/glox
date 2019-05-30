package main

type Stmt interface {
	Accept(StmtVisitor) (interface{}, *RuntimeError)
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) (interface{}, *RuntimeError)
	VisitPrintStmt(*PrintStmt) (interface{}, *RuntimeError)
	VisitVarStmt(*VarStmt) (interface{}, *RuntimeError)
	VisitBlockStmt(*BlockStmt) (interface{}, *RuntimeError)
}

type ExpressionStmt struct {
	Expression Expr
}

func (t *ExpressionStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitExpressionStmt(t)
}

type PrintStmt struct {
	Expression Expr
}

func (t *PrintStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitPrintStmt(t)
}

type VarStmt struct {
	Name *Token
	Initializer Expr
}

func (t *VarStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitVarStmt(t)
}

type BlockStmt struct {
	Statements []Stmt
}

func (t *BlockStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitBlockStmt(t)
}

