package main

// type AstPrinter struct{}
//
// func (a *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, *RuntimeError) {
// 	return a.parenthesize("group", expr.Expression), nil
// }
//
// func (a *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, *RuntimeError) {
// 	return fmt.Sprintf("%v", expr.Value), nil
// }
//
// func (a *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, *RuntimeError) {
// 	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
// }
//
// func (a *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, *RuntimeError) {
// 	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
// }
//
// func (a *AstPrinter) VisitVarExpr(expr *VarExpr) (interface{}, *RuntimeError) {
// 	return "", nil
// }
//
// func (a *AstPrinter) VisitAssignExpr(expr *AssignExpr) (interface{}, *RuntimeError) {
// 	return "", nil
// }
//
// func (a *AstPrinter) Print(expr Expr) {
// 	val, _ := expr.Accept(a)
// 	result := val.(string)
// 	fmt.Printf("%s\n", result)
// }
//
// func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
// 	res := "("
// 	res += name
// 	for _, expr := range exprs {
// 		res += " "
// 		val, _ := expr.Accept(a)
// 		res += val.(string)
// 	}
// 	res += ")"
// 	return res
// }
