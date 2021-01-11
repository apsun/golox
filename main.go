package main

import (
	"bufio"
	"fmt"
	"github.com/apsun/golox/lox"
	"io/ioutil"
	"os"
	"time"
)

func clock(args []lox.Value) (lox.Value, lox.RuntimeException) {
	now := float64(time.Now().UnixNano()) / 1e9
	return lox.NewNumber(now), nil
}

func run(source string, env *lox.Environment, allowExpr bool) bool {
	scanner := lox.NewScanner(source)
	tokens, errs := scanner.ScanTokens()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return false
	}

	env.DefineNative("clock", lox.NewNativeFn(0, "clock", clock))

	parser := lox.NewParser(tokens)
	stmts, errs := parser.ParseStatements()
	if len(errs) > 0 {
		if allowExpr {
			// Try to parse as an expression. If that works, then
			// print out the result of evaluating the expression
			// instead of trying to get a full statement. If neither
			// work, then still show the errors from trying to parse
			// as a statement.
			parser := lox.NewParser(tokens)
			expr, errs := parser.ParseExpression()
			if len(errs) == 0 {
				resolver := lox.NewResolver()
				rerrs := resolver.ResolveExpression(expr)
				if len(rerrs) > 0 {
					for _, err := range rerrs {
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					return false
				}

				value, err := expr.Evaluate(env)
				if err == nil {
					fmt.Printf("%v\n", value.Repr())
				} else {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return false
				}
				return true
			}
		}

		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return false
	}

	resolver := lox.NewResolver()
	rerrs := resolver.ResolveStatements(stmts)
	if len(rerrs) > 0 {
		for _, err := range rerrs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return false
	}

	for _, stmt := range stmts {
		err := stmt.Execute(env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return false
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

	ok := run(string(content), env, false)
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
		run(line, env, true)
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
