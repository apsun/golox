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
	token   *Token
	usages  int
	defined bool
}

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeMethod
	FunctionTypeInitializer
)

type Resolver struct {
	scopes          []map[string]*localVar
	errors          []*ResolverError
	currentFunction FunctionType
}

func NewResolver() *Resolver {
	return &Resolver{
		scopes: []map[string]*localVar{
			map[string]*localVar{},
		},
		errors:          []*ResolverError{},
		currentFunction: FunctionTypeNone,
	}
}

func (r *Resolver) currentScope() map[string]*localVar {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) AddError(token Token, message string) {
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
			r.AddError(*v.token, fmt.Sprintf("'%s' declared but not used", name))
		}
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) Declare(name Token) {
	scope := r.currentScope()
	_, ok := scope[name.lexeme]
	if ok {
		r.AddError(name, fmt.Sprintf("'%s' already declared in this scope", name.lexeme))
	}
	scope[name.lexeme] = &localVar{
		token:   &name,
		usages:  0,
		defined: false,
	}
}

func (r *Resolver) Define(name Token) {
	scope := r.currentScope()
	v, ok := scope[name.lexeme]
	if !ok {
		panic("called Define before Declare")
	}
	v.defined = true
}

func (r *Resolver) DeclareAndDefineNative(name string) {
	scope := r.currentScope()
	_, ok := scope[name]
	if ok {
		panic(fmt.Sprintf("duplicate declaration of '%s'", name))
	}
	scope[name] = &localVar{
		token:   nil,
		usages:  1, // Since we're doing it, suppress unused errors
		defined: true,
	}
}

func (r *Resolver) IsDefined(name Token) bool {
	scope := r.currentScope()
	v, ok := scope[name.lexeme]
	return !ok || v.defined
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

func (r *Resolver) ResolveFunction(e FnExpr, ty FunctionType) {
	oldTy := r.beginFunction(ty)
	defer r.endFunction(oldTy)

	r.BeginScope()
	defer r.EndScope()

	for _, param := range e.parameters {
		r.Declare(param)
		r.Define(param)
	}

	for _, stmt := range e.body {
		stmt.Resolve(r)
	}
}

func (r *Resolver) beginFunction(ty FunctionType) FunctionType {
	old := r.currentFunction
	r.currentFunction = ty
	return old
}

func (r *Resolver) endFunction(prevTy FunctionType) {
	r.currentFunction = prevTy
}

func (r *Resolver) CurrentFunction() FunctionType {
	return r.currentFunction
}
