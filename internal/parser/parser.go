package parser

import (
	"fmt"
	"os"

	"slices"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type ParseError struct {
	Token   token.Token
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Token.Line, e.Token.Lexeme, e.Message)
}

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
		func() {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(ParseError); ok {
						p.HadError = true
						p.synchronize()
					} else {
						panic(r)
					}
				}
			}()

			stmt := p.declaration()
			if stmt != nil {
				statements = append(statements, stmt)
			}
		}()
	}
	return statements
}

// ParseExpressions for backward compatibility
func (p *Parser) ParseExpressions() []ast.Expr {
	var expressions []ast.Expr
	for !p.isAtEnd() {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(ParseError); ok {
						p.HadError = true
						p.synchronize()
					} else {
						panic(r)
					}
				}
			}()

			expr := p.expression()
			if expr != nil {
				expressions = append(expressions, expr)
			}
		}()
	}
	return expressions
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.FUN) {
		return p.function("function")
	}
	return p.statement()
}

func (p *Parser) function(kind string) ast.Stmt {
	name := p.consume(token.IDENTIFIER, fmt.Sprintf("expect %s name", kind))
	p.consume(token.LEFT_PAREN, fmt.Sprintf("expect '(' after %s name", kind))

	parameters := make([]token.Token, 0)
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "can't have more than 255 parameters")
			}

			param := p.consume(token.IDENTIFIER, "expect parameter name")
			parameters = append(parameters, *param)

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	p.consume(token.RIGHT_PAREN, "expect ')' after parameters")
	p.consume(token.LEFT_BRACE, fmt.Sprintf("expect '{' before %s body", kind))
	body := p.block()

	return &ast.FunctionStmt{
		Name:       *name,
		Parameters: parameters,
		Body:       body,
	}
}

func (p *Parser) varDeclaration() ast.Stmt {
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
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return &ast.BlockStmt{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	p.consume(token.SEMICOLON, "expect ';' after return value")
	return &ast.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "expect '(' after 'for'")

	var initializer ast.Stmt = nil
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "expect ';' after loop condition")

	var increment ast.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "expect ')' after for clauses")

	body := p.statement()
	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{body, &ast.ExpressionStmt{Expr: increment}},
		}
	}

	if condition == nil {
		condition = &ast.LiteralExpr{Value: true}
	}
	body = &ast.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.BlockStmt{Statements: []ast.Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "expect '(' after 'while'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "expect ')' after condition")
	body := p.statement()
	return &ast.WhileStmt{Body: body, Condition: condition}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "expect '(' after 'if'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "expect ')' after if condition")

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
		stmt := p.declaration()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	p.consume(token.RIGHT_BRACE, "expect '}' after block")
	return stmts
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "expect ';' after value")
	return &ast.PrintStmt{Expr: expr}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "expect ';' after expression")
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
	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	arguments := make([]ast.Expr, 0)
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "can't have more than 255 arguments")
			}
			arguments = append(arguments, p.expression())

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	parenToken := p.consume(token.RIGHT_PAREN, "expect ')' after arguments")
	return &ast.CallExpr{Callee: callee, Arguments: arguments, Paren: *parenToken}
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
		p.consume(token.RIGHT_PAREN, "expect ')' after expression")
		return &ast.GroupingExpr{Expr: expr}
	}
	if p.match(token.IDENTIFIER) {
		return &ast.VariableExpr{Name: p.previous()}
	}

	p.error(p.peek(), "expect expression")
	return nil
}

func (p *Parser) match(types ...token.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
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
	panic(ParseError{Token: token, Message: message})
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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
