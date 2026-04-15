package expr

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
