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
	var printer ast.AstPrinter

	for {
		if p.isAtEnd() {
			break
		}

		expr := p.ParseFactorExpr()
		if expr == nil {
			p.advance()
			continue
		}

		str, _ := expr.Accept(&printer)
		fmt.Println(str)

		p.advance()
	}
}

func (p *Parser) ParseFactorExpr() ast.Expr {
	expr := p.ParseUnaryExpr()
	p.advance()

	if p.isAtEnd() {
		return expr
	}

	for p.peek().Type == token.STAR || p.peek().Type == token.SLASH {
		op := p.peek()
		p.advance()
		right := p.ParseUnaryExpr()
		p.advance()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) ParseGroupingExpr() ast.Expr {
	p.advance()
	expr := p.ParseFactorExpr()
	p.advance()

	return &ast.GroupingExpr{
		Expr: expr,
	}
}

func (p *Parser) ParseUnaryExpr() ast.Expr {
	if p.peek().Type != token.MINUS && p.peek().Type != token.BANG {
		return p.ParsePrimaryExpr()
	}

	p.advance()
	return &ast.UnaryExpr{
		Operator: p.peekPrevious(),
		Right:    p.ParseUnaryExpr(),
	}
}

func (p *Parser) ParsePrimaryExpr() ast.Expr {
	t := p.peek()

	switch t.Type {
	case token.NUMBER:
		fallthrough
	case token.STRING:
		return &ast.LiteralExpr{Value: t.Literal}

	case token.TRUE:
		fallthrough
	case token.FALSE:
		fallthrough
	case token.NIL:
		return &ast.LiteralExpr{Value: t.Lexeme}

	case token.LEFT_PAREN:
		return p.ParseGroupingExpr()

	default:
		return nil
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

func (p *Parser) peekPrevious() *token.Token {
	return p.tokens[p.current-1]
}
