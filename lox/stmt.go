package lox

import (
	"fmt"
)

type Stmt interface {
	Execute(env *Environment) *RuntimeError
}

type ExprStmt struct {
	expression Expr
}

func (s ExprStmt) Execute(env *Environment) *RuntimeError {
	_, err := s.expression.Evaluate(env)
	return err
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

type VarStmt struct {
	name        Token
	initializer *Expr
}

func (s VarStmt) Execute(env *Environment) *RuntimeError {
	var value Value = NewNil()
	if s.initializer != nil {
		var err *RuntimeError
		value, err = (*s.initializer).Evaluate(env)
		if err != nil {
			return err
		}
	}

	env.Define(s.name.lexeme, value)
	return nil
}
