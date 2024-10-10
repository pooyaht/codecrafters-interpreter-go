package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() {
	for {
		t := p.peek()

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

		case token.NUMBER:
			expr := ast.LiteralExpr{Value: t.Literal}
			str, _ := expr.Accept(&astPrinter)
			fmt.Println(str)

		default:
			return
		}

		p.advance()
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) advance() {
	p.current++
}

func (p *Parser) peek() *token.Token {
	if p.isAtEnd() {
		return nil
	}
	return p.tokens[p.current]
}
