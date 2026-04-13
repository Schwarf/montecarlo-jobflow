package expr

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdentifier

	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenExponent
	TokenLeftParen
	TokenRightParen
	TokenComma
)

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int
}

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenExponent:
		return "^"
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenComma:
		return ","
	default:
		return fmt.Sprintf("TokenType(%d)", int(t))
	}
}
