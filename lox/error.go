package lox

import (
	"fmt"
	"os"
)

var hadError = false

func reportErrorAtToken(token Token, message string) {
	if token.ty == EOF {
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
