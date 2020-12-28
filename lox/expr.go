package lox

import (
	"fmt"
)

type Expr interface {
	Evaluate() (Value, *RuntimeError)
}

// Binary expression

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf(
		"BinaryExpr{left: %v, operator: %v, right: %v}",
		e.left,
		e.operator,
		e.right,
	)
}

func (e BinaryExpr) Evaluate() (Value, *RuntimeError) {
	left, err := e.left.Evaluate()
	if err != nil {
		return nil, err
	}

	right, err := e.right.Evaluate()
	if err != nil {
		return nil, err
	}

	switch e.operator.ty {
	case TokenTypeMinus,
		TokenTypeSlash,
		TokenTypeStar,
		TokenTypeGreater,
		TokenTypeGreaterEqual,
		TokenTypeLess,
		TokenTypeLessEqual:

		l := left.AsNumber()
		r := right.AsNumber()
		if l == nil || r == nil {
			return nil, NewRuntimeError(
				e.operator,
				fmt.Sprintf(
					"%s operands must be numbers",
					e.operator.lexeme,
				),
			)
		}

		switch e.operator.ty {
		case TokenTypeMinus:
			return NewNumber(*l - *r), nil
		case TokenTypeSlash:
			return NewNumber(*l / *r), nil
		case TokenTypeStar:
			return NewNumber(*l * *r), nil
		case TokenTypeGreater:
			return NewBool(*l > *r), nil
		case TokenTypeGreaterEqual:
			return NewBool(*l >= *r), nil
		case TokenTypeLess:
			return NewBool(*l < *r), nil
		case TokenTypeLessEqual:
			return NewBool(*l <= *r), nil
		default:
			panic("unreachable")
		}
	case TokenTypePlus:
		ln := left.AsNumber()
		rn := right.AsNumber()
		if ln != nil && rn != nil {
			return NewNumber(*ln + *rn), nil
		}

		ls := left.AsString()
		rs := right.AsString()
		if ls != nil && rs != nil {
			return NewString(*ls + *rs), nil
		}

		return nil, NewRuntimeError(
			e.operator,
			"+ operands must be numbers or strings",
		)
	case TokenTypeBangEqual:
		return NewBool(!left.Equal(right)), nil
	case TokenTypeEqualEqual:
		return NewBool(left.Equal(right)), nil
	default:
		panic(fmt.Sprintf("unknown binary operator: %v", e.operator.ty))
	}
}

// Grouping expression

type GroupingExpr struct {
	expression Expr
}

func (e GroupingExpr) String() string {
	return fmt.Sprintf("GroupingExpr{expression: %v}", e.expression)
}

func (e GroupingExpr) Evaluate() (Value, *RuntimeError) {
	return e.expression.Evaluate()
}

// Literal expression

type LiteralExpr struct {
	value interface{}
}

func (e LiteralExpr) String() string {
	return fmt.Sprintf("LiteralExpr{value: %#v}", e.value)
}

func (e LiteralExpr) Evaluate() (Value, *RuntimeError) {
	switch v := e.value.(type) {
	case nil:
		return NewNil(), nil
	case bool:
		return NewBool(v), nil
	case float64:
		return NewNumber(v), nil
	case string:
		return NewString(v), nil
	default:
		panic(fmt.Sprintf("unknown literal type: %T", v))
	}
}

// Unary expression

type UnaryExpr struct {
	operator Token
	right    Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf(
		"UnaryExpr{operator: %v, right: %v}",
		e.operator,
		e.right,
	)
}

func (e UnaryExpr) Evaluate() (Value, *RuntimeError) {
	r, err := e.right.Evaluate()
	if err != nil {
		return nil, err
	}

	switch e.operator.ty {
	case TokenTypeBang:
		return NewBool(!r.AsBool()), nil
	case TokenTypeMinus:
		rn := r.AsNumber()
		if rn == nil {
			return nil, NewRuntimeError(
				e.operator,
				"unary - operand must be a number",
			)
		}
		return NewNumber(-*rn), nil
	default:
		panic(fmt.Sprintf("unknown unary operator: %v", e.operator.ty))
	}
}

type TernaryExpr struct {
	cond  Expr
	left  Expr
	right Expr
}

func (e TernaryExpr) String() string {
	return fmt.Sprintf(
		"TernaryExpr{cond: %v, left: %v, right: %v}",
		e.cond,
		e.left,
		e.right,
	)
}

func (e TernaryExpr) Evaluate() (Value, *RuntimeError) {
	cond, err := e.cond.Evaluate()
	if cond != nil {
		return nil, err
	}

	if cond.AsBool() {
		return e.left.Evaluate()
	} else {
		return e.right.Evaluate()
	}
}
