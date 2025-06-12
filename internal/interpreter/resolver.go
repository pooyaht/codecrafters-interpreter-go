package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type functionType int

const (
	NONE functionType = iota
	FUNCTION
	METHOD
	INITIALIZER
)

type Resolver struct {
	interpreter     Interpreter
	scopes          []map[string]bool
	currentFunction functionType
	isInsideClass   bool
}

func NewResolver(interpreter Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NONE,
		isInsideClass:   false,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	r.beginScope()
	defer r.endScope()
	return r.Resolve(stmt.Statements)
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	return r.resolveExpr(stmt.Expr)
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (any, error) {
	if len(r.scopes) != 0 {
		if _, exists := r.scopes[len(r.scopes)-1][stmt.Name.Lexeme]; exists {
			return nil, fmt.Errorf("[Line %d] Error at '%v': Already a function with this name in this scope", stmt.Name.Line, stmt.Name.Lexeme)
		}
	}
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return r.resolveFunction(*stmt, FUNCTION)
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	if _, err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}

	if _, err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}

	if stmt.ElseBranch != nil {
		if _, err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	return r.resolveExpr(stmt.Expr)
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	if r.currentFunction == NONE {
		return nil, fmt.Errorf("[Line %d] Can't return from top-level code", stmt.Keyword.Line)
	}

	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			return nil, fmt.Errorf("[Line %d] Can't return a value from initializer", stmt.Keyword.Line)
		}
		return r.resolveExpr(stmt.Value)
	}

	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	if len(r.scopes) != 0 {
		if _, exists := r.scopes[len(r.scopes)-1][stmt.Name.Lexeme]; exists {
			return nil, fmt.Errorf("[Line %d] Error at '%v': Already a variable with this name in this scope", stmt.Name.Line, stmt.Name.Lexeme)
		}
	}

	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		if _, err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	if _, err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}

	return r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) (any, error) {
	prevIsInsideClass := r.isInsideClass
	r.isInsideClass = true
	defer func() {
		r.isInsideClass = prevIsInsideClass
	}()

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil &&
		stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		return nil, fmt.Errorf("line %d: a class can't inherit from itself", stmt.Superclass.Name.Line)
	}
	if stmt.Superclass != nil {
		r.resolveExpr(stmt.Superclass)
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true
	for _, method := range stmt.Methods {
		functionType := METHOD
		if method.Name.Lexeme == "init" {
			functionType = INITIALIZER
		}
		_, err := r.resolveFunction(method, functionType)
		if err != nil {
			return nil, err
		}
	}
	r.endScope()

	return nil, nil
}

func (r *Resolver) VisitAssignmentExpr(expr *ast.AssignmentExpr) (any, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	return r.resolveLocal(expr, expr.Name)

}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		if _, err := r.resolveExpr(argument); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return r.resolveExpr(expr.Expr)
}

func (r *Resolver) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	return r.resolveExpr(expr.Object)
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	_, err := r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	return r.resolveExpr(expr.Value)
}

func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) (any, error) {
	if len(r.scopes) != 0 {
		if val, exists := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; exists && !val {
			return nil, fmt.Errorf("[Line %d] Error at '%v': Can't read local variable in its own initializer", expr.Name.Line, expr.Name.Lexeme)
		}
	}

	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) VisitThisExpr(expr *ast.ThisExpr) (any, error) {
	if !r.isInsideClass {
		return nil, fmt.Errorf("[Line %d] Can't use 'this' outside of a class", expr.Keyword.Line)
	}
	return r.resolveLocal(expr, expr.Keyword)
}

func (r *Resolver) resolveFunction(stmt ast.FunctionStmt, decleration functionType) (any, error) {
	previousFunctionType := r.currentFunction
	r.currentFunction = decleration
	defer func() {
		r.currentFunction = previousFunctionType
	}()

	r.beginScope()
	defer r.endScope()

	for _, token := range stmt.Parameters {
		if len(r.scopes) != 0 {
			if _, exists := r.scopes[len(r.scopes)-1][token.Lexeme]; exists {
				return nil, fmt.Errorf("[Line %d] Error at '%v': Already a parameter with this name in this scope", stmt.Name.Line, stmt.Name.Lexeme)
			}
		}
		r.declare(token)
		r.define(token)
	}

	return r.Resolve(stmt.Body)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) (any, error) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if _, exists := scope[name.Lexeme]; exists {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			break
		}
	}

	return nil, nil
}

func (r *Resolver) Resolve(statements []ast.Stmt) (any, error) {
	for _, statement := range statements {
		_, err := r.resolveStmt(statement)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) resolveStmt(statement ast.Stmt) (any, error) {
	return statement.Accept(r)
}

func (r *Resolver) resolveExpr(expression ast.Expr) (any, error) {
	return expression.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool, 0))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}
