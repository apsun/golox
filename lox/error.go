package lox

import (
	"fmt"
	"os"
)

var hadError = false

func reportError(line int, message string) {
	fmt.Fprintf(os.Stderr, "error on line %d: %s\n", line, message)
}

func HadError() bool {
	return hadError
}

func ResetError() {
	hadError = false
}
