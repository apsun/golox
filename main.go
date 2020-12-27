package main

import (
	"bufio"
	"fmt"
	"github.com/apsun/golox/lox"
	"io/ioutil"
	"os"
)

func run(source string) bool {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}

	ok := !lox.HadError()
	lox.ResetError()
	return ok
}

func runFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", path)
		os.Exit(1)
	}

	ok := run(string(content))
	if !ok {
		os.Exit(65)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
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