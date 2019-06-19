package main

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeMethod
	FunctionTypeInitializer
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
	ClassTypeSubclass
)

type Resolver struct {
	Interpreter     *Interpreter
	Scopes          []map[string]bool
	CurrentFunction FunctionType
	CurrentClass    ClassType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		Scopes:          []map[string]bool{},
		CurrentFunction: FunctionTypeNone,
		CurrentClass:    ClassTypeNone,
	}
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	newLen := len(r.Scopes) - 1
	r.Scopes = r.Scopes[:newLen]
}

func (r *Resolver) declare(name *Token) {
	if len(r.Scopes) == 0 {
		return
	}

	scope := r.Scopes[len(r.Scopes)-1]

	// Check if exists already in scope and error
	if _, ok := scope[name.Lexeme]; ok {
		lox.errorToken(name, "Variable with this name already declared in this scope.")
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if len(r.Scopes) == 0 {
		return
	}

	scope := r.Scopes[len(r.Scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.Scopes) - 1; i >= 0; i-- {
		scope := r.Scopes[i]
		if _, ok := scope[name.Lexeme]; ok {
			r.Interpreter.resolve(expr, len(r.Scopes)-1-i)
			return
		}
	}

	// Not found. Assume it is global.
}

func (r *Resolver) resolveFunction(function *FunctionStmt, functionType FunctionType) {
	enclosingFunction := r.CurrentFunction
	r.CurrentFunction = functionType
	r.beginScope()

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.resolveStatements(function.Body)
	r.endScope()
	r.CurrentFunction = enclosingFunction
}

// Expressions
func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) (interface{}, *RuntimeError) {
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitVarExpr(expr *VarExpr) (interface{}, *RuntimeError) {
	// Check if defined in local scope and if set to false (ie inside assignment)
	if len(r.Scopes) > 0 {
		localScope := r.Scopes[len(r.Scopes)-1]
		if val, ok := localScope[expr.Name.Lexeme]; ok {
			if val == false {
				lox.errorToken(expr.Name, "Cannot read local variable in its own initializer.")
			}
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *GetExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) (interface{}, *RuntimeError) {
	r.resolveExpression(expr.Value)
	r.resolveExpression(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) (interface{}, *RuntimeError) {
	if r.CurrentClass == ClassTypeNone {
		lox.errorToken(expr.Keyword, "Cannot use 'this' outside of a class.")
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) (interface{}, *RuntimeError) {
	if r.CurrentClass == ClassTypeNone {
		lox.errorToken(expr.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.CurrentClass != ClassTypeSubclass {
		lox.errorToken(expr.Keyword, "Cannot use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

// Statements

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) (interface{}, *RuntimeError) {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) (interface{}, *RuntimeError) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) (interface{}, *RuntimeError) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) (interface{}, *RuntimeError) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionTypeFunction)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) (interface{}, *RuntimeError) {
	if r.CurrentFunction == FunctionTypeNone {
		lox.errorToken(stmt.Keyword, "Cannot return from top-level code.")
	}

	if stmt.Value != nil {
		if r.CurrentFunction == FunctionTypeInitializer {
			lox.errorToken(stmt.Keyword, "Cannot return a value from an initializer.")
		}

		r.resolveExpression(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExpressionStmt) (interface{}, *RuntimeError) {
	r.resolveExpression(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) (interface{}, *RuntimeError) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) (interface{}, *RuntimeError) {
	enclosingClass := r.CurrentClass
	r.CurrentClass = ClassTypeClass
	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil && stmt.Superclass.Name.Lexeme == stmt.Name.Lexeme {
		lox.errorToken(stmt.Superclass.Name, "A class cannot inherit from itself.")
	}

	if stmt.Superclass != nil {
		r.CurrentClass = ClassTypeSubclass
		r.resolveExpression(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.beginScope()
		r.Scopes[len(r.Scopes)-1]["super"] = true
	}

	r.beginScope()
	r.Scopes[len(r.Scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := FunctionTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = FunctionTypeInitializer
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()
	if stmt.Superclass != nil {
		r.endScope()
	}
	r.CurrentClass = enclosingClass
	return nil, nil
}
