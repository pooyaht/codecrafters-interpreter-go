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

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []ast.Expr {
	var expressions []ast.Expr
	var printer ast.AstPrinter
	for !p.isAtEnd() {
		expr := p.expression()
		if expr != nil {
			str, _ := expr.Accept(&printer)
			fmt.Println(str)
			expressions = append(expressions, expr)
		}
	}
	return expressions
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &ast.UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.LiteralExpr{Value: false}
	}
	if p.match(token.TRUE) {
		return &ast.LiteralExpr{Value: true}
	}
	if p.match(token.NIL) {
		return &ast.LiteralExpr{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.match(token.RIGHT_PAREN)
		return &ast.GroupingExpr{Expr: expr}
	}

	return nil
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}
