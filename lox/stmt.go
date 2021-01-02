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
			return err
		}
	}
}
