package expr

import "fmt"

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) current() Token {
	return p.tokens[p.pos]
}

func (p *Parser) isAtEnd() bool {
	return p.current().Type == TokenEOF
}

func (p *Parser) advance() Token {
	current := p.tokens[p.pos]
	if !p.isAtEnd() {
		p.pos++
	}
	return current
}

func (p *Parser) Parse() (Expression, error) {
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.isAtEnd() {
		token := p.current()
		return nil, fmt.Errorf("unexpected token %q at position %d", token.Lexeme, token.Pos)
	}
	return expression, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
	token := p.current()

	switch token.Type {
	case TokenNumber:
		p.advance()
		return &NumberExpression{Value: token.Lexeme}, nil

	case TokenIdentifier:
		p.advance()
		return &VariableExpression{Name: token.Lexeme}, nil

	case TokenLeftParen:
		p.advance()

		expression, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		if p.current().Type != TokenRightParen {
			return nil, fmt.Errorf("expected ')' at position %d but got '%s'", p.current().Pos, p.current().Lexeme)
		}
		p.advance()
		return expression, nil

	default:
		return nil, fmt.Errorf("expected primary expression at position %d, got %q", token.Pos, token.Lexeme)
	}
}
