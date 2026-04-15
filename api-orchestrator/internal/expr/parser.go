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
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		token := p.current()
		if token.Type != TokenPlus && token.Type != TokenMinus {
			break
		}

		operator := token.Type
		p.advance()

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// power is right-associative: 3^4^5 = 3^(4^5)
func (p *Parser) parsePower() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	if p.current().Type == TokenPower {
		operator := p.current().Type
		p.advance()

		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}, nil
	}
	return left, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}

	for {
		token := p.current()
		if token.Type != TokenMultiply && token.Type != TokenDivide {
			break
		}
		operator := token.Type
		p.advance()

		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	token := p.current()

	if token.Type == TokenPlus || token.Type == TokenMinus {
		operator := token.Type
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
	token := p.current()

	switch token.Type {
	case TokenNumber:
		p.advance()
		return &NumberExpression{Value: token.Lexeme}, nil

	case TokenIdentifier:
		if p.peek().Type == TokenLeftParen {
			functionName := token.Lexeme
			p.advance() // consume identifier
			p.advance() // consume '('

			argument, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			if p.current().Type != TokenRightParen {
				return nil, fmt.Errorf("expected ')' after function argument at position %d but got '%s'", p.current().Pos, p.current().Lexeme)
			}
			p.advance() // consume ')'

			return &FunctionCallExpression{
				Name:     functionName,
				Argument: argument,
			}, nil
		}

		p.advance()
		return &VariableExpression{Name: token.Lexeme}, nil

	case TokenLeftParen:
		p.advance()

		expression, err := p.parseExpression()
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

func (p *Parser) peek() Token {
	if p.pos+1 >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos+1]
}
