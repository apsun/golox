package lox

import (
	"fmt"
)

type Stmt interface {
	Execute(env *Environment) *RuntimeError
	Resolve(r *Resolver)
}

type ExprStmt struct {
	expression Expr
}

func (s ExprStmt) Execute(env *Environment) *RuntimeError {
	_, err := s.expression.Evaluate(env)
	return err
}

func (s ExprStmt) Resolve(r *Resolver) {
	s.expression.Resolve(r)
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) Execute(env *Environment) *RuntimeError {
	val, err := s.expression.Evaluate(env)
	if err != nil {
		return err
	}

	fmt.Println(val.String())
	return nil
}

func (s PrintStmt) Resolve(r *Resolver) {
	s.expression.Resolve(r)
}

type VarStmt struct {
	name        Token
	initializer *Expr
}

func (s VarStmt) Execute(env *Environment) *RuntimeError {
	if s.initializer == nil {
		env.Declare(s.name)
	} else {
		value, err := (*s.initializer).Evaluate(env)
		if err != nil {
			return err
		}
		env.Define(s.name, value)
	}
	return nil
}

func (s VarStmt) Resolve(r *Resolver) {
	r.Declare(s.name)
	if s.initializer != nil {
		(*s.initializer).Resolve(r)
	}
	r.Define(s.name)
}

type BlockStmt struct {
	statements []Stmt
}

func (s BlockStmt) Execute(env *Environment) *RuntimeError {
	innerEnv := NewEnvironment(env)
	for _, stmt := range s.statements {
		err := stmt.Execute(innerEnv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s BlockStmt) Resolve(r *Resolver) {
	r.BeginScope()
	defer r.EndScope()

	for _, stmt := range s.statements {
		stmt.Resolve(r)
	}
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch *Stmt
}

func (s IfStmt) Execute(env *Environment) *RuntimeError {
	val, err := s.condition.Evaluate(env)
	if err != nil {
		return err
	}

	if val.Bool() {
		return s.thenBranch.Execute(env)
	} else if s.elseBranch != nil {
		return (*s.elseBranch).Execute(env)
	} else {
		return nil
	}
}

func (s IfStmt) Resolve(r *Resolver) {
	s.condition.Resolve(r)
	s.thenBranch.Resolve(r)
	if s.elseBranch != nil {
		(*s.elseBranch).Resolve(r)
	}
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) Execute(env *Environment) *RuntimeError {
	for {
		val, err := s.condition.Evaluate(env)
		if err != nil {
			return err
		}

		if !val.Bool() {
			return nil
		}

		err = s.body.Execute(env)
		if err != nil {
			if err == breakError {
				return nil
			}
			return err
		}
	}
}

func (s WhileStmt) Resolve(r *Resolver) {
	s.condition.Resolve(r)
	s.body.Resolve(r)
}

type BreakStmt struct{}

// Singleton "error" that we use to break out of loops
var breakError = NewRuntimeErrorNoToken("break")

func (s BreakStmt) Execute(env *Environment) *RuntimeError {
	return breakError
}

func (s BreakStmt) Resolve(r *Resolver) {
	// No-op
}

type FnStmt struct {
	name       Token
	parameters []Token
	body       []Stmt
}

func (s FnStmt) Execute(env *Environment) *RuntimeError {
	fn := NewLoxFn(s)
	env.Define(s.name, fn)
	return nil
}

func (s FnStmt) Resolve(r *Resolver) {
	r.Declare(s.name)
	r.Define(s.name)

	r.BeginScope()
	defer r.EndScope()

	for _, param := range s.parameters {
		r.Declare(param)
		r.Define(param)
	}

	for _, stmt := range s.body {
		stmt.Resolve(r)
	}
}
