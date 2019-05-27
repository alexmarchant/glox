package main

import (
	"strconv"
)

// Scanner is in charge of breaking source string into tokens
type Scanner struct {
	Source string
	Tokens []*Token
	Start int
	Current int
	Line int
}

func makeScanner(source string) *Scanner {
	return &Scanner{
		Source: source,
		Start: 0,
		Current: 0,
		Line: 1,
	}
}

func (s *Scanner) scanTokens() []*Token {
	for !s.isAtEnd() {
		s.Start = s.Current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, &Token{
		Type: EOF,
		Line: s.Line,
	})
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
		case '(':
			s.addToken(LeftParen)
		case ')':
			s.addToken(RightParen)
		case '{':
			s.addToken(LeftBrace)
		case '}':
			s.addToken(RightBrace)
		case ',':
			s.addToken(Comma)
		case '.':
			s.addToken(Dot)
		case '-':
			s.addToken(Minus)
		case '+':
			s.addToken(Plus)
		case ';':
			s.addToken(Semicolon)
		case '*':
			s.addToken(Star)
		case '!':
			if s.match('=') {
				s.addToken(BangEqual)
			} else {
				s.addToken(Bang)
			}
		case '=':
			if s.match('=') {
				s.addToken(EqualEqual)
			} else {
				s.addToken(Equal)
			}
		case '<':
			if s.match('=') {
				s.addToken(LessEqual)
			} else {
				s.addToken(Less)
			}
		case '>':
			if s.match('=') {
				s.addToken(GreaterEqual)
			} else {
				s.addToken(Greater)
			}
		case '/':
			if s.match('/') {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				s.addToken(Slash)
			}
		case ' ', '\r', '\t':
		case '\n':
			s.Line++
		case '"':
			s.string()
		default:
			if isDigit(char) {
				s.number()
			} else if isAlpha(char) {
				s.identifier()
			} else {
				lox.errorLine(s.Line, "Unexpected character.")
			}
	}
}

func (s *Scanner) advance() rune {
	s.Current++
	return rune(s.Source[s.Current - 1])
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenValue(tokenType, nil)
}

func (s *Scanner) addTokenValue(tokenType TokenType, literal LiteralValue) {
	text := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, &Token{
		Type: tokenType,
		Lexeme: text,
		Literal: literal,
		Line: s.Line,
	})
}

func (s *Scanner) match(char rune) bool {
	if s.isAtEnd() {
		return false
	}
	nextChar := rune(s.Source[s.Current])
	if nextChar != char {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) peekNext() rune {
	if s.Current + 1 >= len(s.Source) {
		return rune(0)
	}
	return rune(s.Source[s.Current+1])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		lox.errorLine(s.Line, "Unterminated string.")
	}

	s.advance()

	value := s.Source[s.Start+1:s.Current-1]
	s.addTokenValue(String, &LiteralValueString{Value: value})
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	if err != nil {
		lox.errorLine(s.Line, "Ivalid number.")
	}
	s.addTokenValue(Number, &LiteralValueNumber{Value: value})
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]
	val, ok := keywords[text]
	if ok {
		s.addToken(val)
	} else {
		s.addToken(Identifier)
	}
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func isAlphaNumeric(char rune) bool {
	return isDigit(char) || isAlpha(char)
}