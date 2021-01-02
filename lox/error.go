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

type RuntimeError struct {
	token   *Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   &token,
		message: message,
	}
}

func NewRuntimeErrorNoToken(message string) *RuntimeError {
	return &RuntimeError{
		token:   nil,
		message: message,
	}
}

func (e *RuntimeError) Error() string {
	line := -1
	if e.token != nil {
		line = e.token.line
	}

	return fmt.Sprintf(
		"runtime error on line %d: %s",
		line,
		e.message,
	)
}
