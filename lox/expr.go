package lox

import (
	"fmt"
)

type Expr interface {
	Evaluate(env *Environment) (Value, RuntimeException)
	Resolve(r *Resolver)
}

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

func (e BinaryExpr) Evaluate(env *Environment) (Value, RuntimeException) {
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

func (e BinaryExpr) Resolve(r *Resolver) {
	e.left.Resolve(r)
	e.right.Resolve(r)
}

type GroupingExpr struct {
	expression Expr
}

func (e GroupingExpr) String() string {
	return fmt.Sprintf("GroupingExpr{expression: %v}", e.expression)
}

func (e GroupingExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	return e.expression.Evaluate(env)
}

func (e GroupingExpr) Resolve(r *Resolver) {
	e.expression.Resolve(r)
}

type LiteralExpr struct {
	value interface{}
}

func (e LiteralExpr) String() string {
	return fmt.Sprintf("LiteralExpr{value: %#v}", e.value)
}

func (e LiteralExpr) Evaluate(env *Environment) (Value, RuntimeException) {
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

func (e LiteralExpr) Resolve(r *Resolver) {
	// No-op
}

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

func (e UnaryExpr) Evaluate(env *Environment) (Value, RuntimeException) {
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

func (e UnaryExpr) Resolve(r *Resolver) {
	e.right.Resolve(r)
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

func (e TernaryExpr) Evaluate(env *Environment) (Value, RuntimeException) {
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

func (e TernaryExpr) Resolve(r *Resolver) {
	e.cond.Resolve(r)
	e.left.Resolve(r)
	e.right.Resolve(r)
}

type VariableExpr struct {
	name     Token
	distance *int
}

func (e VariableExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	return env.Get(*e.distance, e.name)
}

func (e VariableExpr) Resolve(r *Resolver) {
	r.CheckDefined(e.name)
	distance := r.ResolveLocal(e.name)
	*e.distance = distance
}

type AssignExpr struct {
	name     Token
	value    Expr
	distance *int
}

func (e AssignExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	value, err := e.value.Evaluate(env)
	if err != nil {
		return nil, err
	}

	err = env.Assign(*e.distance, e.name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (e AssignExpr) Resolve(r *Resolver) {
	e.value.Resolve(r)
	distance := r.ResolveLocal(e.name)
	*e.distance = distance
}

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (e LogicalExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	left, err := e.left.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.operator.ty {
	case TokenTypeAnd:
		if left.Bool() {
			return e.right.Evaluate(env)
		} else {
			return left, nil
		}
	case TokenTypeOr:
		if left.Bool() {
			return left, nil
		} else {
			return e.right.Evaluate(env)
		}
	default:
		panic(fmt.Sprintf("unknown logical operator: %v", e.operator.ty))
	}
}

func (e LogicalExpr) Resolve(r *Resolver) {
	e.left.Resolve(r)
	e.right.Resolve(r)
}

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (e CallExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	callee, err := e.callee.Evaluate(env)
	if err != nil {
		return nil, err
	}

	args := make([]Value, len(e.arguments))
	for i, argExpr := range e.arguments {
		arg, err := argExpr.Evaluate(env)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	fn, ok := callee.(CallableValue)
	if !ok {
		return nil, NewRuntimeError(e.paren, "value is not callable")
	}

	if fn.Arity() != len(args) {
		return nil, NewRuntimeError(
			e.paren,
			fmt.Sprintf(
				"expected %d arguments but got %d",
				fn.Arity(),
				len(args),
			),
		)
	}

	return fn.Call(env, args)
}

func (e CallExpr) Resolve(r *Resolver) {
	e.callee.Resolve(r)
	for _, arg := range e.arguments {
		arg.Resolve(r)
	}
}
