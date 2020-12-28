package lox

import (
	"fmt"
)

type Environment struct {
	values map[string]Value
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]Value{},
	}
}

func (e *Environment) Define(name Token, value Value) *RuntimeError {
	e.values[name.lexeme] = value
	return nil
}

func (e *Environment) Assign(name Token, value Value) *RuntimeError {
	_, ok := e.values[name.lexeme]
	if !ok {
		return NewRuntimeError(
			name,
			fmt.Sprintf("undefined variable %s", name.lexeme),
		)
	}
	e.values[name.lexeme] = value
	return nil
}

func (e *Environment) Get(name Token) (Value, *RuntimeError) {
	value, ok := e.values[name.lexeme]
	if !ok {
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("undefined variable %s", name.lexeme),
		)
	}
	return value, nil
}
