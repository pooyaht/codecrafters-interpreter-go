package parser

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Parser struct {
	tokens   []token.Token
	current  int
	HadError bool
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		current:  0,
		HadError: false,
	}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		stmt := p.decleration()
		if p.HadError {
			return nil
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements
}

// for tests backward compatibility
func (p *Parser) ParseExpressions() []ast.Expr {
	var expressions []ast.Expr
	for !p.isAtEnd() {
		expr := p.expression()
		if p.HadError {
			return nil
		}
		if expr != nil {
			expressions = append(expressions, expr)
		}
	}
	return expressions
}

func (p *Parser) decleration() ast.Stmt {
	if p.match(token.VAR) {
		return p.varDecleration()
	}

	return p.statement()
}

func (p *Parser) varDecleration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "expect variable name")
	var initializer ast.Expr = nil
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "expect ';' after variable declaration")
	return &ast.VarStmt{Name: *name, Initializer: initializer}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return &ast.BlockStmt{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "expect '(' after 'while'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "expect ')' after 'while' condition")
	body := p.statement()
	return &ast.WhileStmt{Body: body, Condition: condition}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "expect '(' after 'if'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "expect ')' after 'if' condition")

	thenBranchStmt := p.statement()
	var elseBranchStmt ast.Stmt = nil
	if p.match(token.ELSE) {
		elseBranchStmt = p.statement()
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranchStmt,
		ElseBranch: elseBranchStmt,
	}
}

func (p *Parser) block() []ast.Stmt {
	stmts := make([]ast.Stmt, 0)
	for !p.isAtEnd() && !p.check(token.RIGHT_BRACE) {
		stmts = append(stmts, p.decleration())
	}

	p.consume(token.RIGHT_BRACE, "expect '}' after block")
	return stmts
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &ast.PrintStmt{Expr: expr}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &ast.ExpressionStmt{Expr: expr}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(*ast.VariableExpr); ok {
			return &ast.AssignmentExpr{
				Name:  varExpr.Name,
				Value: value,
			}
		}

		p.error(equals, "invalid assignment target")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &ast.LogicalExpr{
			Operator: operator,
			Right:    right,
			Left:     expr,
		}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &ast.LogicalExpr{
			Operator: operator,
			Right:    right,
			Left:     expr,
		}
	}

	return expr
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
		right := p.term()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
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
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		return &ast.GroupingExpr{Expr: expr}
	}
	if p.match(token.IDENTIFIER) {
		return &ast.VariableExpr{Name: p.previous()}
	}
	p.error(p.peek(), "Expect expression")
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

func (p *Parser) consume(t token.TokenType, message string) *token.Token {
	if p.check(t) {
		token := p.advance()
		return &token
	}
	p.error(p.peek(), message)
	return nil
}

func (p *Parser) error(token token.Token, message string) {
	p.HadError = true
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, message)
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}
