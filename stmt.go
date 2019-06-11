package main

type Stmt interface {
	Accept(StmtVisitor) (interface{}, *RuntimeError)
}

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt) (interface{}, *RuntimeError)
	VisitIfStmt(*IfStmt) (interface{}, *RuntimeError)
	VisitWhileStmt(*WhileStmt) (interface{}, *RuntimeError)
	VisitFunctionStmt(*FunctionStmt) (interface{}, *RuntimeError)
	VisitReturnStmt(*ReturnStmt) (interface{}, *RuntimeError)
	VisitExpressionStmt(*ExpressionStmt) (interface{}, *RuntimeError)
	VisitVarStmt(*VarStmt) (interface{}, *RuntimeError)
}

type ExpressionStmt struct {
	Expression Expr
}

func (t *ExpressionStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitExpressionStmt(t)
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

type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (t *IfStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitIfStmt(t)
}

type WhileStmt struct {
	Condition Expr
	Body Stmt
}

func (t *WhileStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitWhileStmt(t)
}

type FunctionStmt struct {
	Name *Token
	Params []*Token
	Body []Stmt
}

func (t *FunctionStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitFunctionStmt(t)
}

type ReturnStmt struct {
	Keyword *Token
	Value Expr
}

func (t *ReturnStmt) Accept(visitor StmtVisitor) (interface{}, *RuntimeError) {
	return visitor.VisitReturnStmt(t)
}

