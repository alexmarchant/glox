package main

import (
	"errors"
)

type Parser struct {
	Tokens []*Token
	Current int
}

func (p *Parser) parse() Expr {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	return expr
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
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
			Left: expr,
			Operator: operator,
			Right: right,
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
			Left: expr,
			Operator: operator,
			Right: right,
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
			Left: expr,
			Operator: operator,
			Right: right,
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
			Left: expr,
			Operator: operator,
			Right: right,
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
			Right: right,
		}, nil
	}

	return p.primary()
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