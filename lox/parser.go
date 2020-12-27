package lox

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

var unwindToken = struct{}{}

func (p *Parser) Parse() (ret Expr) {
	defer func() {
		r := recover()
		if r != nil {
			if r == unwindToken {
				ret = nil
			} else {
				panic(r)
			}
		}
	}()
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.comma()
}

func (p *Parser) comma() Expr {
	expr := p.equality()

	for p.match(Comma) {
		operator := p.previous()
		right := p.equality()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr

}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{
			operator: operator,
			right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(False) {
		return LiteralExpr{value: false}
	}

	if p.match(True) {
		return LiteralExpr{value: true}
	}

	if p.match(Nil) {
		return LiteralExpr{value: nil}
	}

	if p.match(Number, String) {
		return LiteralExpr{value: p.previous().literal}
	}

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "expected ')' after expression")
		return GroupingExpr{expression: expr}
	}

	p.unwindWithError(p.peek(), "expected expression")
	panic("unreachable")
}

func (p *Parser) consume(ty TokenType, message string) Token {
	if p.check(ty) {
		return p.advance()
	}

	p.unwindWithError(p.peek(), message)
	panic("unreachable")
}

func (p *Parser) unwindWithError(t Token, message string) {
	reportErrorAtToken(t, message)
	panic(unwindToken)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().ty == Semicolon {
			return
		}

		switch p.peek().ty {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(tys ...TokenType) bool {
	for _, ty := range tys {
		if p.check(ty) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(ty TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ty == ty
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().ty == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
