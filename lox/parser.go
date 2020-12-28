package lox

type Parser struct {
	tokens  []Token
	current int
	errors  []*SyntaxError
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  []*SyntaxError{},
	}
}

var unwindToken = struct{}{}

func (p *Parser) Parse() (ret Expr, errs []*SyntaxError) {
	defer func() {
		r := recover()
		if r != nil {
			if r == unwindToken {
				ret = nil
				errs = p.errors
			} else {
				panic(r)
			}
		}
	}()
	return p.expression(), nil
}

func (p *Parser) expression() Expr {
	return p.comma()
}

func (p *Parser) comma() Expr {
	expr := p.ternary()

	for p.match(TokenTypeComma) {
		operator := p.previous()
		right := p.ternary()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) ternary() Expr {
	expr := p.equality()

	if p.match(TokenTypeQuestion) {
		left := p.expression()
		p.consume(TokenTypeColon, "expected ':' after expression")
		right := p.ternary()
		return TernaryExpr{
			cond:  expr,
			left:  left,
			right: right,
		}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(
		TokenTypeBangEqual,
		TokenTypeEqualEqual,
	) {
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

	for p.match(
		TokenTypeGreater,
		TokenTypeGreaterEqual,
		TokenTypeLess,
		TokenTypeLessEqual,
	) {
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

	for p.match(TokenTypeMinus, TokenTypePlus) {
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

	for p.match(TokenTypeSlash, TokenTypeStar) {
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
	if p.match(TokenTypeBang, TokenTypeMinus) {
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
	if p.match(TokenTypeFalse) {
		return LiteralExpr{value: false}
	}

	if p.match(TokenTypeTrue) {
		return LiteralExpr{value: true}
	}

	if p.match(TokenTypeNil) {
		return LiteralExpr{value: nil}
	}

	if p.match(TokenTypeNumber, TokenTypeString) {
		return LiteralExpr{value: p.previous().literal}
	}

	if p.match(TokenTypeLeftParen) {
		expr := p.expression()
		p.consume(TokenTypeRightParen, "expected ')' after expression")
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
	p.errors = append(p.errors, NewSyntaxError(t.line, &t, message))
	panic(unwindToken)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().ty == TokenTypeSemicolon {
			return
		}

		switch p.peek().ty {
		case TokenTypeClass,
			TokenTypeFun,
			TokenTypeVar,
			TokenTypeFor,
			TokenTypeIf,
			TokenTypeWhile,
			TokenTypePrint,
			TokenTypeReturn:
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
	return p.peek().ty == TokenTypeEOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
