package lox

import (
	"fmt"
)

type SyntaxError struct {
	line    int
	token   *Token
	message string
}

func NewSyntaxError(line int, token *Token, message string) *SyntaxError {
	return &SyntaxError{
		line:    line,
		token:   token,
		message: message,
	}
}

func (e *SyntaxError) Error() string {
	if e.token != nil {
		return fmt.Sprintf(
			"syntax error on line %d at %s: %s",
			e.line,
			e.token.lexeme,
			e.message,
		)
	} else {
		return fmt.Sprintf(
			"syntax error on line %d: %s",
			e.line,
			e.message,
		)
	}
}

type RuntimeException interface{}

type RuntimeError struct {
	token   Token
	message string
}

type BreakException struct{}

type ReturnException struct {
	value Value
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf(
		"runtime error on line %d: %s",
		e.token.line,
		e.message,
	)
}
