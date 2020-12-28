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

func (e *Environment) Declare(name Token) {
	e.values[name.lexeme] = nil
}

func (e *Environment) Define(name Token, value Value) {
	e.values[name.lexeme] = &value
}

func (e *Environment) Assign(name Token, value Value) *RuntimeError {
	_, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}

		return NewRuntimeError(
			name,
			fmt.Sprintf("undefined variable %s", name.lexeme),
		)
	}

	e.values[name.lexeme] = &value
	return nil
}

func (e *Environment) Get(name Token) (Value, *RuntimeError) {
	value, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}

		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("undefined variable %s", name.lexeme),
		)
	}

	if value == nil {
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("using uninitialized variable %s", name.lexeme),
		)
	}

	return *value, nil
}
