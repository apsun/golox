package lox

import (
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]*Value
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		enclosing: outer,
		values:    map[string]*Value{},
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	curr := e
	for distance != 0 {
		if curr.enclosing == nil {
			return curr
		}
		curr = curr.enclosing
		distance--
	}
	return curr
}

func (e *Environment) Declare(name Token) {
	e.values[name.lexeme] = nil
}

func (e *Environment) Define(name Token, value Value) {
	e.values[name.lexeme] = &value
}

func (e *Environment) DefineNative(name string, value Value) {
	e.values[name] = &value
}

func (e *Environment) Assign(distance int, name Token, value Value) RuntimeException {
	e.ancestor(distance).values[name.lexeme] = &value
	return nil
}

func (e *Environment) Get(distance int, name Token) (Value, RuntimeException) {
	value, ok := e.ancestor(distance).values[name.lexeme]
	if !ok {
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("using undeclared variable '%s'", name.lexeme),
		)
	}
	if value == nil {
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("using uninitialized variable '%s'", name.lexeme),
		)
	}
	return *value, nil
}

func (e *Environment) GetNative(distance int, name string) Value {
	value, ok := e.ancestor(distance).values[name]
	if !ok {
		panic(fmt.Sprintf("using undeclared variable '%s'", name))
	}
	if value == nil {
		panic(fmt.Sprintf("using uninitialized variable '%s'", name))
	}
	return *value
}
