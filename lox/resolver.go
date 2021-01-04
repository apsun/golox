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

type localVar struct {
	token   Token
	usages  int
	defined bool
}

type Resolver struct {
	scopes []map[string]*localVar
	errors []*ResolverError
}

func NewResolver() *Resolver {
	return &Resolver{
		scopes: []map[string]*localVar{
			map[string]*localVar{},
		},
		errors: []*ResolverError{},
	}
}

func (r *Resolver) currentScope() map[string]*localVar {
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
	r.scopes = append(r.scopes, map[string]*localVar{})
}

func (r *Resolver) EndScope() {
	scope := r.scopes[len(r.scopes)-1]
	for name, v := range scope {
		if v.usages == 0 && name[0] != '_' {
			r.addError(v.token, fmt.Sprintf("'%s' declared but not used", name))
		}
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) Declare(name Token) {
	scope := r.currentScope()
	_, ok := scope[name.lexeme]
	if ok {
		r.addError(name, fmt.Sprintf("'%s' already declared in this scope", name.lexeme))
	}
	scope[name.lexeme] = &localVar{
		token:   name,
		usages:  0,
		defined: false,
	}
}

func (r *Resolver) Define(name Token) {
	scope := r.currentScope()
	scope[name.lexeme].defined = true
}

func (r *Resolver) CheckDefined(name Token) {
	scope := r.currentScope()
	v, ok := scope[name.lexeme]
	if ok && !v.defined {
		r.addError(name, fmt.Sprintf("cannot refer to '%s' in its own initializer", name.lexeme))
	}
}

func (r *Resolver) ResolveLocal(name Token) int {
	for i := range r.scopes {
		scope := r.scopes[len(r.scopes)-1-i]
		v, ok := scope[name.lexeme]
		if ok {
			v.usages++
			return i
		}
	}
	return -1
}
