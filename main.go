package main

import (
	"bufio"
	"fmt"
	"github.com/apsun/golox/lox"
	"io/ioutil"
	"os"
)

func run(source string, env *lox.Environment) bool {
	scanner := lox.NewScanner(source)
	tokens, errs := scanner.ScanTokens()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return false
	}

	parser := lox.NewParser(tokens)
	exprs, errs := parser.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return false
	}

	for _, expr := range exprs {
		err := expr.Execute(env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}

	return true
}

func runFile(path string) {
	env := lox.NewEnvironment(nil)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", path)
		os.Exit(1)
	}

	ok := run(string(content), env)
	if !ok {
		os.Exit(65)
	}
}

func runPrompt() {
	env := lox.NewEnvironment(nil)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line, env)
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read stdin failed: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [file]\n", os.Args[0])
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
