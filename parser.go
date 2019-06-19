package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	Tokens  []*Token
	Current int
}

func (p *Parser) parse() []Stmt {
	statements := []Stmt{}

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() Stmt {
	var err error
	var statement Stmt
	if p.match(Class) {
		statement, err = p.classDeclaration()
	} else if p.match(Fun) {
		statement, err = p.function("function")
	} else if p.match(Var) {
		statement, err = p.varDeclaration()
	} else {
		statement, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil
	}
	return statement
}

func (p *Parser) varDeclaration() (Stmt, error) {
	var initializer Expr
	var err error

	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(LeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &BlockStmt{
			Statements: statements,
		}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Initializer
	var initializer Stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Condition
	var condition Expr
	if !p.check(Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(Semicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Increment
	var increment Expr
	if !p.check(RightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	// Body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Add increment to body
	if increment != nil {
		body = &BlockStmt{
			Statements: []Stmt{
				body,
				&ExpressionStmt{
					Expression: increment,
				},
			},
		}
	}

	// Add condition to while loop w/ body inside
	if condition == nil {
		condition = &LiteralExpr{
			Value: true,
		}
	}
	body = &WhileStmt{
		Condition: condition,
		Body:      body,
	}

	// Add initializer before while loop
	if initializer != nil {
		body = &BlockStmt{
			Statements: []Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	// Condition
	_, err := p.consume(LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	// Then
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Else
	var elseBranch Stmt
	if p.match(Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()

	var value Expr
	var err error

	if !p.check(Semicolon) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	// Condition
	_, err := p.consume(LeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	// Body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := []Stmt{}

	for !p.check(RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, err := p.consume(RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ExpressionStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) classDeclaration() (*ClassStmt, error) {
	name, err := p.consume(Identifier, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *VarExpr
	if p.match(Less) {
		_, err = p.consume(Identifier, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &VarExpr{
			Name: p.previous(),
		}
	}

	_, err = p.consume(LeftBrace, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []*FunctionStmt{}

	for !p.check(RightBrace) && !p.isAtEnd() {
		method, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	_, err = p.consume(RightBrace, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &ClassStmt{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}

func (p *Parser) function(kind string) (*FunctionStmt, error) {
	// Func name
	name, err := p.consume(Identifier, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	// Open paren
	_, err = p.consume(LeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	// Params
	parameters := []*Token{}
	if !p.check(RightParen) {
		for {
			if len(parameters) >= 8 {
				_ = p.error(p.peek(), "Cannot have more than 8 parameters.")
			}

			newParam, err := p.consume(Identifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, newParam)

			if !p.match(Comma) {
				break
			}
		}
	}

	// Close paren
	_, err = p.consume(RightParen, fmt.Sprintf("Expect ')' after parameters."))
	if err != nil {
		return nil, err
	}

	// Left brace
	_, err = p.consume(LeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	// Body
	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*VarExpr); ok {
			return &AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		} else if getExpr, ok := expr.(*GetExpr); ok {
			return &SetExpr{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}, nil
		}

		err = p.error(equals, "Invalid assignment target.")
		return nil, err
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) addition() (Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) multiplication() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(Dot) {
			name, err := p.consume(Identifier, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &GetExpr{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(RightParen) {
		for {
			if len(arguments) > 8 {
				_ = p.error(p.peek(), "Cannot have more than 8 arguments.")
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)
			if !p.match(Comma) {
				break
			}
		}
	}

	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primary() (Expr, error) {
	switch {
	case p.match(False):
		return &LiteralExpr{
			Value: false,
		}, nil
	case p.match(True):
		return &LiteralExpr{
			Value: true,
		}, nil
	case p.match(Nil):
		return &LiteralExpr{
			Value: nil,
		}, nil
	case p.match(Number, String):
		return &LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	case p.match(Identifier):
		return &VarExpr{
			Name: p.previous(),
		}, nil
	case p.match(Super):
		keyword := p.previous()
		_, err := p.consume(Dot, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(Identifier, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return &SuperExpr{
			Keyword: keyword,
			Method:  method,
		}, nil
	case p.match(This):
		return &ThisExpr{
			Keyword: p.previous(),
		}, nil
	case p.match(LeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &GroupingExpr{
			Expression: expr,
		}, nil
	default:
		err := p.error(p.peek(), "Exprected expression.")
		return nil, err
	}
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	err := p.error(p.peek(), message)
	return nil, err
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() *Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}

		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}

func (p *Parser) error(token *Token, msg string) error {
	lox.errorToken(token, msg)
	return errors.New(msg)
}
