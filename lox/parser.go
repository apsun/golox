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

func (p *Parser) ParseStatements() ([]Stmt, []*SyntaxError) {
	stmts := []Stmt{}
	for !p.isAtEnd() {
		decl := p.parseStatement()
		if decl != nil {
			stmts = append(stmts, *decl)
		}
	}
	return stmts, p.errors
}

func (p *Parser) ParseExpression() (Expr, []*SyntaxError) {
	expr := p.parseExpression()
	if expr == nil {
		return nil, p.errors
	}
	if !p.isAtEnd() {
		p.addError(p.peek(), "trailing junk")
	}
	return *expr, p.errors
}

func (p *Parser) parseStatement() (ret *Stmt) {
	defer func() {
		r := recover()
		if r != nil {
			if r == unwindToken {
				p.synchronize()
				ret = nil
			} else {
				panic(r)
			}
		}
	}()

	stmt := p.declaration()
	return &stmt
}

func (p *Parser) parseExpression() (ret *Expr) {
	defer func() {
		r := recover()
		if r != nil {
			if r == unwindToken {
				p.synchronize()
				ret = nil
			} else {
				panic(r)
			}
		}
	}()

	expr := p.expression()
	return &expr
}

func (p *Parser) declaration() Stmt {
	if p.match(TokenTypeVar) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(TokenTypeIdentifier, "expected variable name")

	var initializer *Expr = nil
	if p.match(TokenTypeEqual) {
		expr := p.expression()
		initializer = &expr
	}

	p.consume(TokenTypeSemicolon, "expected ';' after variable declaration")
	return VarStmt{
		name:        name,
		initializer: initializer,
	}
}

func (p *Parser) statement() Stmt {
	if p.match(TokenTypeFor) {
		return p.forStatement()
	}
	if p.match(TokenTypeIf) {
		return p.ifStatement()
	}
	if p.match(TokenTypePrint) {
		return p.printStatement()
	}
	if p.match(TokenTypeWhile) {
		return p.whileStatement()
	}
	if p.match(TokenTypeLeftBrace) {
		return p.block()
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(TokenTypeLeftParen, "expected '(' after 'if'")

	var initializer *Stmt
	if p.match(TokenTypeSemicolon) {
		initializer = nil
	} else if p.match(TokenTypeVar) {
		tmp := p.varDeclaration()
		initializer = &tmp
	} else {
		tmp := p.expressionStatement()
		initializer = &tmp
	}

	var condition *Expr = nil
	if !p.check(TokenTypeSemicolon) {
		tmp := p.expression()
		condition = &tmp
	}
	p.consume(TokenTypeSemicolon, "expected ';' after condition")

	var increment *Expr = nil
	if !p.check(TokenTypeRightParen) {
		tmp := p.expression()
		increment = &tmp
	}
	p.consume(TokenTypeRightParen, "expected ')' after for clauses")

	body := p.statement()

	if increment != nil {
		body = BlockStmt{
			statements: []Stmt{
				body,
				ExprStmt{expression: *increment},
			},
		}
	}

	if condition == nil {
		body = WhileStmt{condition: LiteralExpr{value: true}, body: body}
	} else {
		body = WhileStmt{condition: *condition, body: body}
	}

	if initializer != nil {
		body = BlockStmt{
			statements: []Stmt{
				*initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(TokenTypeLeftParen, "expected '(' after 'if'")
	condition := p.expression()
	p.consume(TokenTypeRightParen, "expected ')' after if condition")

	thenBranch := p.statement()
	var elseBranch *Stmt = nil
	if p.match(TokenTypeElse) {
		tmp := p.statement()
		elseBranch = &tmp
	}

	return IfStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(TokenTypeLeftParen, "expected '(' after 'while'")
	condition := p.expression()
	p.consume(TokenTypeRightParen, "expected ')' after while condition")
	body := p.statement()

	return WhileStmt{
		condition: condition,
		body:      body,
	}
}

func (p *Parser) block() Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() && !p.check(TokenTypeRightBrace) {
		decl := p.parseStatement()
		if decl != nil {
			statements = append(statements, *decl)
		}
	}
	stmt := BlockStmt{
		statements: statements,
	}
	p.consume(TokenTypeRightBrace, "expected '}' after block")
	return stmt
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(TokenTypeSemicolon, "expected ';' after value")
	return PrintStmt{
		expression: value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(TokenTypeSemicolon, "expected ';' after expression")
	return ExprStmt{
		expression: expr,
	}
}

func (p *Parser) expression() Expr {
	return p.comma()
}

func (p *Parser) comma() Expr {
	expr := p.assignment()

	for p.match(TokenTypeComma) {
		operator := p.previous()
		right := p.assignment()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) assignment() Expr {
	expr := p.ternary()

	if p.match(TokenTypeEqual) {
		equals := p.previous()
		value := p.assignment()

		// Left hand side needs to be a variable expression
		// for now. We reach into it and grab the name of the
		// variable.
		varExpr, ok := expr.(VariableExpr)
		if ok {
			name := varExpr.name
			return AssignExpr{
				name:  name,
				value: value,
			}
		}

		p.addError(equals, "invalid assignment target")
	}

	return expr
}

func (p *Parser) ternary() Expr {
	expr := p.or()

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

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(TokenTypeOr) {
		op := p.previous()
		right := p.and()
		expr = LogicalExpr{
			left:     expr,
			operator: op,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(TokenTypeAnd) {
		op := p.previous()
		right := p.equality()
		expr = LogicalExpr{
			left:     expr,
			operator: op,
			right:    right,
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

	if p.match(TokenTypeIdentifier) {
		return VariableExpr{name: p.previous()}
	}

	if p.match(TokenTypeLeftParen) {
		expr := p.expression()
		p.consume(TokenTypeRightParen, "expected ')' after expression")
		return GroupingExpr{expression: expr}
	}

	p.addError(p.peek(), "expected expression")
	panic(unwindToken)
}

func (p *Parser) consume(ty TokenType, message string) Token {
	if p.check(ty) {
		return p.advance()
	}

	p.addError(p.peek(), message)
	panic(unwindToken)
}

func (p *Parser) addError(t Token, message string) {
	p.errors = append(p.errors, NewSyntaxError(t.line, &t, message))
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
