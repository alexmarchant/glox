package main

import (
	"fmt"
)

type AstPrinter struct {
	Output string
}

func (a *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) {
	a.addParenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) {
	a.Output += expr.Value.String()
}

func (a *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) {
	a.addParenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) {
	a.addParenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) Print(expr Expr) {
	expr.Accept(a)
	fmt.Printf("%s\n", a.Output)
}

func (a *AstPrinter) addParenthesize(name string, exprs ...Expr) {
	a.Output += "("
	a.Output += name
	for _, expr := range exprs {
		a.Output += " "
		expr.Accept(a)
	}
	a.Output += ")"
}