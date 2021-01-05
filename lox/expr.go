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

		if left.Type() != TypeNumber || right.Type() != TypeNumber {
			return nil, NewRuntimeError(
				e.operator,
				fmt.Sprintf(
					"%s operands must be numbers",
					e.operator.lexeme,
				),
			)
		}

		l := left.(Number).Float()
		r := right.(Number).Float()

		switch e.operator.ty {
		case TokenTypeMinus:
			return NewNumber(l - r), nil
		case TokenTypeSlash:
			if r == 0 {
				return nil, NewRuntimeError(
					e.operator,
					fmt.Sprintf("division by zero"),
				)
			}
			return NewNumber(l / r), nil
		case TokenTypeStar:
			return NewNumber(l * r), nil
		case TokenTypeGreater:
			return NewBool(l > r), nil
		case TokenTypeGreaterEqual:
			return NewBool(l >= r), nil
		case TokenTypeLess:
			return NewBool(l < r), nil
		case TokenTypeLessEqual:
			return NewBool(l <= r), nil
		default:
			panic("unreachable")
		}
	case TokenTypePlus:
		if left.Type() == TypeNumber && right.Type() == TypeNumber {
			ln := left.(Number).Float()
			rn := right.(Number).Float()
			return NewNumber(ln + rn), nil
		}

		if left.Type() == TypeString || right.Type() == TypeString {
			ls := left.String()
			rs := right.String()
			return NewString(ls + rs), nil
		}

		return nil, NewRuntimeError(
			e.operator,
			"+ operands must be numbers or strings",
		)
	case TokenTypeBangEqual:
		return NewBool(!left.Equal(right)), nil
	case TokenTypeEqualEqual:
		return NewBool(left.Equal(right)), nil
	case TokenTypeComma:
		return right, nil
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
		if r.Type() != TypeNumber {
			return nil, NewRuntimeError(
				e.operator,
				"unary - operand must be a number",
			)
		}
		rn := r.(Number).Float()
		return NewNumber(-rn), nil
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
	if !r.IsDefined(e.name) {
		r.AddError(
			e.name,
			fmt.Sprintf(
				"cannot refer to '%s' in its own initializer",
				e.name.lexeme,
			),
		)
	}
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

func (e CallExpr) callLoxFn(fn *LoxFn, args []Value) (Value, RuntimeException) {
	declaration, env := fn.FnWithEnv()

	calleeEnv := NewEnvironment(env)
	for i, arg := range args {
		name := declaration.parameters[i]
		calleeEnv.Define(name, arg)
	}

	var result Value = NewNil()

	for _, stmt := range declaration.body {
		err := stmt.Execute(calleeEnv)
		if err != nil {
			ret, ok := err.(ReturnException)
			if ok {
				result = ret.value
				break
			}
			return nil, err
		}
	}

	// Initializers act as if they return the instance
	if fn.IsInit() {
		return env.GetNative(0, "this"), nil
	}

	return result, nil
}

func (e CallExpr) newInstance(class *Class, args []Value) (Value, RuntimeException) {
	instance := NewInstance(class)
	initializer := class.Initializer()
	if initializer != nil {
		return e.callLoxFn(instance.Bind(*initializer), args)
	}
	return instance, nil
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

	callable, ok := callee.(Callable)
	if !ok {
		return nil, NewRuntimeError(e.paren, "value is not callable")
	}

	arity := callable.Arity()
	if arity != len(args) {
		return nil, NewRuntimeError(
			e.paren,
			fmt.Sprintf(
				"expected %d arguments but got %d",
				arity,
				len(args),
			),
		)
	}

	switch callable := callable.(type) {
	case *NativeFn:
		return callable.Fn()(args)
	case *LoxFn:
		return e.callLoxFn(callable, args)
	case *Class:
		return e.newInstance(callable, args)
	default:
		panic("unreachable")
	}
}

func (e CallExpr) Resolve(r *Resolver) {
	e.callee.Resolve(r)
	for _, arg := range e.arguments {
		arg.Resolve(r)
	}
}

type FnExpr struct {
	parameters []Token
	body       []Stmt
}

func (e FnExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	return NewLoxFn(nil, e, env, false), nil
}

func (e FnExpr) Resolve(r *Resolver) {
	r.ResolveFunction(e, FunctionTypeFunction)
}

type GetExpr struct {
	object Expr
	name   Token
}

func (e GetExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	object, err := e.object.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if object.Type() != TypeInstance {
		return nil, NewRuntimeError(e.name, "only instances have properties")
	}

	return object.(*Instance).Get(e.name)
}

func (e GetExpr) Resolve(r *Resolver) {
	e.object.Resolve(r)
}

type SetExpr struct {
	object Expr
	name   Token
	value  Expr
}

func (e SetExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	object, err := e.object.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if object.Type() != TypeInstance {
		return nil, NewRuntimeError(e.name, "only instances have properties")
	}

	value, err := e.value.Evaluate(env)
	if err != nil {
		return nil, err
	}

	object.(*Instance).Set(e.name, value)
	return value, nil
}

func (e SetExpr) Resolve(r *Resolver) {
	e.value.Resolve(r)
	e.object.Resolve(r)
}

type ThisExpr struct {
	keyword  Token
	distance *int
}

func (e ThisExpr) Evaluate(env *Environment) (Value, RuntimeException) {
	return env.Get(*e.distance, e.keyword)
}

func (e ThisExpr) Resolve(r *Resolver) {
	ty := r.CurrentFunction()
	if ty != FunctionTypeMethod && ty != FunctionTypeInitializer {
		r.AddError(e.keyword, "cannot use this outside method")
	}
	*e.distance = r.ResolveLocal(e.keyword)
}
