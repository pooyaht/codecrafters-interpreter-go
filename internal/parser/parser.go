package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Parser struct {
	scanner scanner.Scanner
}

func NewParser(scanner scanner.Scanner) Parser {
	return Parser{
		scanner: scanner,
	}
}

func (p *Parser) Parse() {
	for {
		t, err := p.scanner.Scan()

		if err != nil {
			break
		}
		if t == nil {
			continue
		}

		astPrinter := ast.AstPrinter{}

		switch t.Type {
		case token.TRUE:
			expr := ast.LiteralExpr{Value: true}
			str, _ := expr.Accept(&astPrinter)
			fmt.Println(str)

		case token.FALSE:
			expr := ast.LiteralExpr{Value: false}
			str, _ := expr.Accept(&astPrinter)
			fmt.Println(str)

		case token.NIL:
			expr := ast.LiteralExpr{Value: nil}
			str, _ := expr.Accept(&astPrinter)
			fmt.Println(str)

		default:
			return
		}

	}
}
