package lox

import (
	"fmt"
)

type Expr interface {
	Evaluate(env *Environment) (Value, *RuntimeError)
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

func (e BinaryExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	left, err := e.left.Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := e.right.Evaluate(env)
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

		l := left.CastNumber()
		r := right.CastNumber()
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
			if *r == 0 {
				return nil, NewRuntimeError(
					e.operator,
					fmt.Sprintf("division by zero"),
				)
			}
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
		ln := left.CastNumber()
		rn := right.CastNumber()
		if ln != nil && rn != nil {
			return NewNumber(*ln + *rn), nil
		}

		ls := left.CastString()
		rs := right.CastString()
		if ls != nil && rs != nil {
			return NewString(*ls + *rs), nil
		}

		if ls != nil {
			return NewString(*ls + right.String()), nil
		}

		if rs != nil {
			return NewString(left.String() + *rs), nil
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

func (e GroupingExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	return e.expression.Evaluate(env)
}

// Literal expression

type LiteralExpr struct {
	value interface{}
}

func (e LiteralExpr) String() string {
	return fmt.Sprintf("LiteralExpr{value: %#v}", e.value)
}

func (e LiteralExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
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

func (e UnaryExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	r, err := e.right.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.operator.ty {
	case TokenTypeBang:
		return NewBool(!r.Bool()), nil
	case TokenTypeMinus:
		rn := r.CastNumber()
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

func (e TernaryExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	cond, err := e.cond.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if cond.Bool() {
		return e.left.Evaluate(env)
	} else {
		return e.right.Evaluate(env)
	}
}

type VariableExpr struct {
	name Token
}

func (e VariableExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	return env.Get(e.name)
}

type AssignExpr struct {
	name  Token
	value Expr
}

func (e AssignExpr) Evaluate(env *Environment) (Value, *RuntimeError) {
	value, err := e.value.Evaluate(env)
	if err != nil {
		return nil, err
	}

	err = env.Assign(e.name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
