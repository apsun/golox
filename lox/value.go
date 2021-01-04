package lox

import (
	"fmt"
	"strconv"
)

type Type int

const (
	TypeNil Type = iota
	TypeBool
	TypeNumber
	TypeString
	TypeCallable
)

type Value interface {
	Type() Type
	Bool() bool
	Equal(other Value) bool
	String() string
	Repr() string
}

// nil
type Nil struct{}

var nilInstance = Nil{}

func NewNil() Nil {
	return nilInstance
}

func (x Nil) Type() Type {
	return TypeNil
}

func (x Nil) Bool() bool {
	return false
}

func (x Nil) Equal(other Value) bool {
	return x == other
}

func (x Nil) String() string {
	return "nil"
}

func (x Nil) Repr() string {
	return x.String()
}

// bool
type Bool struct{ value bool }

func NewBool(value bool) Bool {
	return Bool{value: value}
}

func (x Bool) Type() Type {
	return TypeBool
}

func (x Bool) Bool() bool {
	return x.value
}

func (x Bool) Equal(other Value) bool {
	return other.Type() == TypeBool && x.value == other.(Bool).value
}

func (x Bool) String() string {
	return strconv.FormatBool(x.value)
}

func (x Bool) Repr() string {
	return x.String()
}

// number
type Number struct{ value float64 }

func NewNumber(value float64) Number {
	return Number{value: value}
}

func (x Number) Type() Type {
	return TypeNumber
}

func (x Number) Bool() bool {
	return true
}

func (x Number) Equal(other Value) bool {
	return other.Type() == TypeNumber && x.value == other.(Number).value
}

func (x Number) String() string {
	return strconv.FormatFloat(x.value, 'f', -1, 64)
}

func (x Number) Repr() string {
	return x.String()
}

func (x Number) Float() float64 {
	return x.value
}

// string
type String struct{ value string }

func NewString(value string) String {
	return String{value: value}
}

func (x String) Type() Type {
	return TypeString
}

func (x String) Bool() bool {
	return true
}

func (x String) Equal(other Value) bool {
	return other.Type() == TypeString && x.value == other.(String).value
}

func (x String) String() string {
	return x.value
}

func (x String) Repr() string {
	return fmt.Sprintf("%q", x.value)
}

// native fn
type NativeFnPtr func(args []Value) (Value, RuntimeException)

type NativeFn struct {
	arity int
	name  string
	fn    NativeFnPtr
}

func NewNativeFn(arity int, name string, fn NativeFnPtr) *NativeFn {
	return &NativeFn{
		arity: arity,
		name:  name,
		fn:    fn,
	}
}

func (x *NativeFn) Type() Type {
	return TypeCallable
}

func (x *NativeFn) Bool() bool {
	return true
}

func (x *NativeFn) Equal(other Value) bool {
	return x == other
}

func (x *NativeFn) String() string {
	return fmt.Sprintf("<native fn '%s'>", x.name)
}

func (x *NativeFn) Repr() string {
	return x.String()
}

func (x *NativeFn) Arity() int {
	return x.arity
}

func (x *NativeFn) Fn() NativeFnPtr {
	return x.fn
}

// lox fn
type LoxFn struct {
	name        *string
	declaration FnExpr
	env         *Environment
}

func NewLoxFn(name *string, declaration FnExpr, env *Environment) *LoxFn {
	return &LoxFn{
		name:        name,
		declaration: declaration,
		env:         env,
	}
}

func (x *LoxFn) Type() Type {
	return TypeCallable
}

func (x *LoxFn) Bool() bool {
	return true
}

func (x *LoxFn) Equal(other Value) bool {
	return x == other
}

func (x *LoxFn) String() string {
	if x.name != nil {
		return fmt.Sprintf("<fn '%s'>", *x.name)
	} else {
		return "<anonymous fn>"
	}
}

func (x *LoxFn) Repr() string {
	return x.String()
}

func (x *LoxFn) Arity() int {
	return len(x.declaration.parameters)
}

func (x *LoxFn) FnWithEnv() (FnExpr, *Environment) {
	return x.declaration, x.env
}
