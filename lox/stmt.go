package lox

import (
	"fmt"
)

type Stmt interface {
	Execute(env *Environment) RuntimeException
	Resolve(r *Resolver)
}

type ExprStmt struct {
	expression Expr
}

func (s ExprStmt) Execute(env *Environment) RuntimeException {
	_, err := s.expression.Evaluate(env)
	return err
}

func (s ExprStmt) Resolve(r *Resolver) {
	s.expression.Resolve(r)
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) Execute(env *Environment) RuntimeException {
	val, err := s.expression.Evaluate(env)
	if err != nil {
		return err
	}

	fmt.Println(val.String())
	return nil
}

func (s PrintStmt) Resolve(r *Resolver) {
	s.expression.Resolve(r)
}

type VarStmt struct {
	name        Token
	initializer *Expr
}

func (s VarStmt) Execute(env *Environment) RuntimeException {
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

func (s VarStmt) Resolve(r *Resolver) {
	r.Declare(s.name)
	if s.initializer != nil {
		(*s.initializer).Resolve(r)
	}
	r.Define(s.name)
}

type BlockStmt struct {
	statements []Stmt
}

func (s BlockStmt) Execute(env *Environment) RuntimeException {
	innerEnv := NewEnvironment(env)
	for _, stmt := range s.statements {
		err := stmt.Execute(innerEnv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s BlockStmt) Resolve(r *Resolver) {
	r.BeginScope()
	defer r.EndScope()

	for _, stmt := range s.statements {
		stmt.Resolve(r)
	}
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch *Stmt
}

func (s IfStmt) Execute(env *Environment) RuntimeException {
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

func (s IfStmt) Resolve(r *Resolver) {
	s.condition.Resolve(r)
	s.thenBranch.Resolve(r)
	if s.elseBranch != nil {
		(*s.elseBranch).Resolve(r)
	}
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) Execute(env *Environment) RuntimeException {
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
			_, ok := err.(BreakException)
			if ok {
				return nil
			}
			return err
		}
	}
}

func (s WhileStmt) Resolve(r *Resolver) {
	s.condition.Resolve(r)
	s.body.Resolve(r)
}

type BreakStmt struct{}

func (s BreakStmt) Execute(env *Environment) RuntimeException {
	return BreakException{}
}

func (s BreakStmt) Resolve(r *Resolver) {
	// No-op
}

type FnStmt struct {
	name     Token
	function FnExpr
}

func (s FnStmt) Execute(env *Environment) RuntimeException {
	name := s.name.lexeme
	fn := NewLoxFn(&name, s.function, env, false, false)
	env.Define(s.name, fn)
	return nil
}

func (s FnStmt) Resolve(r *Resolver) {
	r.Declare(s.name)
	r.Define(s.name)
	s.function.Resolve(r)
}

type MethodStmt struct {
	FnStmt
	isProperty bool
}

func (s MethodStmt) Execute(env *Environment) RuntimeException {
	panic("should not be called")
}

func (s MethodStmt) Resolve(r *Resolver) {
	panic("should not be called")
}

type ReturnStmt struct {
	keyword Token
	value   *Expr
}

func (s ReturnStmt) Execute(env *Environment) RuntimeException {
	if s.value == nil {
		return ReturnException{value: NewNil()}
	}

	value, err := (*s.value).Evaluate(env)
	if err != nil {
		return err
	}
	return ReturnException{value: value}
}

func (s ReturnStmt) Resolve(r *Resolver) {
	ty := r.CurrentFunction()
	if ty == FunctionTypeNone {
		r.AddError(s.keyword, "cannot return outside function")
	} else if ty == FunctionTypeInitializer && s.value != nil {
		r.AddError(s.keyword, "cannot return value from initializer")
	}

	if s.value != nil {
		(*s.value).Resolve(r)
	}
}

type ClassStmt struct {
	name         Token
	methods      []MethodStmt
	classMethods []MethodStmt
}

func (s ClassStmt) Execute(env *Environment) RuntimeException {
	env.Declare(s.name)

	classMethods := map[string]*LoxFn{}
	for _, method := range s.classMethods {
		name := method.name.lexeme
		fn := NewLoxFn(&name, method.function, env, false, method.isProperty)
		classMethods[method.name.lexeme] = fn
	}
	metaclass := NewClass(nil, s.name.lexeme+" metaclass", classMethods)

	methods := map[string]*LoxFn{}
	for _, method := range s.methods {
		name := method.name.lexeme
		isInit := (name == "init")
		fn := NewLoxFn(&name, method.function, env, isInit, method.isProperty)
		methods[method.name.lexeme] = fn
	}
	class := NewClass(metaclass, s.name.lexeme, methods)

	env.Define(s.name, class)
	return nil
}

func (s ClassStmt) Resolve(r *Resolver) {
	r.Declare(s.name)
	r.Define(s.name)

	r.BeginScope()
	defer r.EndScope()

	r.DeclareAndDefineNative("this")

	for _, method := range s.classMethods {
		r.ResolveFunction(method.function, FunctionTypeMethod)
	}

	for _, method := range s.methods {
		ty := FunctionTypeMethod
		if method.name.lexeme == "init" {
			ty = FunctionTypeInitializer
		}
		r.ResolveFunction(method.function, ty)
	}
}
