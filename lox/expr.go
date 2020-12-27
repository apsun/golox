package lox

import (
	"fmt"
)

type Expr interface{}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr{left: %v, operator: %v, right: %v}", e.left, e.operator, e.right)
}

type GroupingExpr struct {
	expression Expr
}

func (e GroupingExpr) String() string {
	return fmt.Sprintf("GroupingExpr{expression: %v}", e.expression)
}

type LiteralExpr struct {
	value interface{}
}

func (e LiteralExpr) String() string {
	return fmt.Sprintf("LiteralExpr{value: %#v}", e.value)
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("UnaryExpr{operator: %v, right: %v}", e.operator, e.right)
}
