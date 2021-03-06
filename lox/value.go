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
	TypeFn
	TypeClass
	TypeInstance
)

type Value interface {
	Type() Type
	Bool() bool
	Equal(other Value) bool
	String() string
	Repr() string
}

type Callable interface {
	Value
	Arity() int
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
	return TypeFn
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
	isInit      bool
	isProperty  bool
}

func NewLoxFn(
	name *string,
	declaration FnExpr,
	env *Environment,
	isInit bool,
	isProperty bool,
) *LoxFn {
	return &LoxFn{
		name:        name,
		declaration: declaration,
		env:         env,
		isInit:      isInit,
		isProperty:  isProperty,
	}
}

func (x *LoxFn) Type() Type {
	return TypeFn
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

func (x *LoxFn) IsInit() bool {
	return x.isInit
}

func (x *LoxFn) IsProperty() bool {
	return x.isProperty
}

// common interface for classes and instances
type Fielder interface {
	Get(name Token) (Value, RuntimeException)
	Set(name Token, value Value)
}

// class
type Class struct {
	Instance
	name       string
	superclass **Class
	methods    map[string]*LoxFn
}

func NewClass(
	metaclass **Class,
	name string,
	superclass **Class,
	methods map[string]*LoxFn,
) *Class {
	return &Class{
		Instance: Instance{
			class:  metaclass,
			fields: map[string]Value{},
		},
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (x *Class) Type() Type {
	return TypeClass
}

func (x *Class) Bool() bool {
	return true
}

func (x *Class) Equal(other Value) bool {
	return x == other
}

func (x *Class) String() string {
	return fmt.Sprintf("<class '%s'>", x.name)
}

func (x *Class) Repr() string {
	return x.String()
}

func (x *Class) method(name string) **LoxFn {
	method, ok := x.methods[name]
	if ok {
		return &method
	}
	if x.superclass != nil {
		return (*x.superclass).method(name)
	}
	return nil
}

func (x *Class) initializer() **LoxFn {
	return x.method("init")
}

func (x *Class) Arity() int {
	initializer := x.initializer()
	if initializer != nil {
		return (*initializer).Arity()
	}
	return 0
}

// class instance
type Instance struct {
	class  **Class
	fields map[string]Value
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  &class,
		fields: map[string]Value{},
	}
}

func (x *Instance) Type() Type {
	return TypeInstance
}

func (x *Instance) Bool() bool {
	return true
}

func (x *Instance) Equal(other Value) bool {
	return x == other
}

func (x *Instance) String() string {
	return fmt.Sprintf("<instance of class '%s'>", x.Class().name)
}

func (x *Instance) Repr() string {
	return x.String()
}

func (x *Instance) Class() *Class {
	return *x.class
}

func (x *Instance) bind(method *LoxFn) *LoxFn {
	env := NewEnvironment(method.env)
	env.DefineNative("this", x)
	return NewLoxFn(method.name, method.declaration, env, method.isInit, method.isProperty)
}

func (x *Instance) Initializer() **LoxFn {
	init := x.Class().initializer()
	if init == nil {
		return nil
	}
	tmp := x.bind(*init)
	return &tmp
}

func (x *Instance) MethodAtClass(name Token, class *Class) (*LoxFn, RuntimeException) {
	method := class.method(name.lexeme)
	if method != nil {
		return x.bind(*method), nil
	}

	return nil, NewRuntimeError(
		name,
		fmt.Sprintf("undefined method '%s'", class.name, name.lexeme),
	)
}

func (x *Instance) Get(name Token) (Value, RuntimeException) {
	value, ok := x.fields[name.lexeme]
	if ok {
		return value, nil
	}

	method := x.Class().method(name.lexeme)
	if method != nil {
		return x.bind(*method), nil
	}

	return nil, NewRuntimeError(
		name,
		fmt.Sprintf("undefined method/property/field '%s'", name.lexeme),
	)
}

func (x *Instance) Set(name Token, value Value) {
	x.fields[name.lexeme] = value
}
