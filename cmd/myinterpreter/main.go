package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "tokenize" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(string(fileContents))
		hadError := false
		for {
			t, err := scanner.Scan()
			if err != nil {
				hadError = true
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			if t == nil {
				continue
			}

			fmt.Println(t)

			if t.Type == token.EOF {
				break
			}
		}

		if hadError {
			os.Exit(65)
		}
	} else if command == "parse" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(string(fileContents))
		tokens, err := scanner.ScanTokens()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		parser := parser.NewParser(tokens)
		astPrinter := &ast.AstPrinter{}
		nodes := parser.ParseExpressions()
		if parser.HadError {
			os.Exit(65)
		}

		for _, node := range nodes {
			str, _ := node.Accept(astPrinter)
			fmt.Println(str)
		}

	} else if command == "evaluate" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(string(fileContents))
		tokens, err := scanner.ScanTokens()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		parser := parser.NewParser(tokens)
		nodes := parser.ParseExpressions()
		if parser.HadError {
			os.Exit(65)
		}

		interpreter := interpreter.Interpreter{}
		for _, node := range nodes {
			val, err := node.Accept(&interpreter)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(70)
			} else if val == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(val)
			}
		}
	} else if command == "run" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(string(fileContents))
		tokens, err := scanner.ScanTokens()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		parser := parser.NewParser(tokens)
		nodes := parser.Parse()
		if parser.HadError {
			os.Exit(65)
		}

		interpreter := interpreter.Interpreter{}
		for _, node := range nodes {
			val, err := node.Accept(&interpreter)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(70)
			} else if val == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(val)
			}
		}
	}
}
