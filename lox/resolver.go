package lox

import (
	"fmt"
)

type ResolverError struct {
	token   Token
	message string
}

func newResolverError(token Token, message string) *ResolverError {
	return &ResolverError{
		token:   token,
		message: message,
	}
}

func (e *ResolverError) String() string {
	return fmt.Sprintf(
		"resolver error on line %d at %s: %s",
		e.token.line,
		e.token.lexeme,
		e.message,
	)
}

type Resolver struct {
	scopes []map[string]bool
	errors []*ResolverError
}

func NewResolver() *Resolver {
	return &Resolver{
		scopes: []map[string]bool{
			map[string]bool{},
		},
		errors: []*ResolverError{},
	}
}

func (r *Resolver) currentScope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) addError(token Token, message string) {
	r.errors = append(r.errors, newResolverError(token, message))
}

func (r *Resolver) ResolveStatements(stmts []Stmt) []*ResolverError {
	for _, stmt := range stmts {
		stmt.Resolve(r)
	}
	return r.errors
}

func (r *Resolver) ResolveExpression(expr Expr) []*ResolverError {
	expr.Resolve(r)
	return r.errors
}

func (r *Resolver) BeginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) EndScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) Declare(name Token) {
	scope := r.currentScope()
	_, ok := scope[name.lexeme]
	if ok {
		r.addError(name, fmt.Sprintf("'%s' already declared in this scope", name.lexeme))
	}
	scope[name.lexeme] = false
}

func (r *Resolver) Define(name Token) {
	scope := r.currentScope()
	scope[name.lexeme] = true
}

func (r *Resolver) CheckDefined(name Token) {
	scope := r.currentScope()
	ready, ok := scope[name.lexeme]
	if !ready && ok {
		r.addError(name, fmt.Sprintf("cannot refer to '%s' in its own initializer", name.lexeme))
	}
}

func (r *Resolver) ResolveLocal(name Token) int {
	for i := range r.scopes {
		scope := r.scopes[len(r.scopes)-1-i]
		_, ok := scope[name.lexeme]
		if ok {
			return i
		}
	}
	return -1
}
