package lox

import (
	"fmt"
)

type Stmt interface {
	Execute() *RuntimeError
}

type ExprStmt struct {
	expression Expr
}

func (s ExprStmt) Execute() *RuntimeError {
	_, err := s.expression.Evaluate()
	return err
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) Execute() *RuntimeError {
	val, err := s.expression.Evaluate()
	if err != nil {
		return err
	}

	fmt.Println(val.String())
	return nil
}
