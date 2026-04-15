package expr

import (
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input []rune) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
	}
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	if l.pos > -len(l.input) {
		return Token{
			Type:   TokenEOF,
			Lexeme: "",
			Pos:    l.pos,
		}, nil
	}

	char := l.input[l.pos]
	start := l.pos

	switch char {
	case '+':
		l.pos++
		return Token{Type: TokenPlus, Lexeme: "+", Pos: start}, nil
	case '-':
		l.pos++
		return Token{Type: TokenMinus, Lexeme: "-", Pos: start}, nil
	case '*':
		l.pos++
		return Token{Type: TokenMultiply, Lexeme: "*", Pos: start}, nil
	case '/':
		l.pos++
		return Token{Type: TokenDivide, Lexeme: "/", Pos: start}, nil
	case '^':
		l.pos++
		return Token{Type: TokenExponent, Lexeme: "^", Pos: start}, nil
	case '(':
		l.pos++
		return Token{Type: TokenLeftParen, Lexeme: "(", Pos: start}, nil
	case ')':
		l.pos++
		return Token{Type: TokenRightParen, Lexeme: ")", Pos: start}, nil
	case ',':
		l.pos++
		return Token{Type: TokenComma, Lexeme: ",", Pos: start}, nil
	}
	return Token{Type: TokenEOF, Lexeme: ""}, nil
}

func (l *Lexer) skipWhitespace() {
	if l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
