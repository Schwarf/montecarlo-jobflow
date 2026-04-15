package expr

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
	}
}

func LexAll(input string) ([]Token, error) {
	lexer := NewLexer(input)
	var tokens []Token

	for {
		tok, err := lexer.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}

	return tokens, nil
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
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
		return Token{Type: TokenPower, Lexeme: "^", Pos: start}, nil
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

	if isDigit(char) || char == '.' {
		return l.readNumber()
	}

	if isIdentifierStart(char) {
		return l.readIdentifier(), nil
	}

	return Token{}, fmt.Errorf("unexpected character %q at position %d", char, start)
}

func (l *Lexer) readNumber() (Token, error) {
	start := l.pos
	seenDot := false

	if l.input[l.pos] == '.' {
		seenDot = true
		l.pos++
		if l.pos >= len(l.input) || !isDigit(l.input[l.pos]) {
			return Token{}, fmt.Errorf("invalid number at position %d", start)
		}
	}

	for l.pos < len(l.input) {
		char := l.input[l.pos]

		if isDigit(char) {
			l.pos++
			continue
		}

		if char == '.' {
			if seenDot {
				break
			}
			seenDot = true
			l.pos++
			continue
		}
		break
	}

	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		expPos := l.pos
		l.pos++

		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}

		if l.pos >= len(l.input) || !isDigit(l.input[l.pos]) {
			return Token{}, fmt.Errorf("invalid scientific notation at position %d", expPos)
		}

		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	return Token{Type: TokenNumber, Lexeme: string(l.input[start:l.pos]), Pos: start}, nil

}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	l.pos++

	for l.pos < len(l.input) && isIdentifierPart(l.input[l.pos]) {
		l.pos++
	}

	return Token{
		Type:   TokenIdentifier,
		Lexeme: string(l.input[start:l.pos]),
		Pos:    start,
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
