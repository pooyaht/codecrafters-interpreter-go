package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := scanner.NewScanner(string(fileContents))

	var errorOccurred bool

	for {
		t, err := scanner.Scan()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			errorOccurred = true
			continue
		}
		fmt.Println(t)

		if t.Type == token.EOF {
			break
		}
	}

	if errorOccurred {
		os.Exit(65)
	}
}
