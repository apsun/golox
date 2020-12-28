package lox

import (
	"fmt"
	"os"
)

var hadError = false

type RuntimeError struct {
	token   Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

func (e *RuntimeError) Error() string {
	return e.message
}

func reportErrorAtToken(token Token, message string) {
	if token.ty == TokenTypeEOF {
		reportErrorAt(token.line, "end", message)
	} else {
		reportErrorAt(token.line, token.lexeme, message)
	}
}

func reportErrorAt(line int, where string, message string) {
	fmt.Fprintf(
		os.Stderr,
		"error on line %d at %s: %s\n",
		line, where, message,
	)
	hadError = true
}

func reportError(line int, message string) {
	fmt.Fprintf(os.Stderr, "error on line %d: %s\n", line, message)
	hadError = true
}

func HadError() bool {
	return hadError
}

func ResetError() {
	hadError = false
}
