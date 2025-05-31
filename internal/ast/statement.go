package ast

import (
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type PrintStmt struct {
	Expr Expr
}

func (s *PrintStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrintStmt(s)
}

type ExpressionStmt struct {
	Expr Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionStmt(s)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (s *VarStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitVarStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitBlockStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s *IfStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitWhileStmt(s)
}

type FunctionStmt struct {
	Name       token.Token
	Parameters []token.Token
	Body       []Stmt
}

func (s *FunctionStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitFunctionStmt(s)
}

type ReturnStmt struct {
	Value   Expr
	Keyword token.Token
}

func (s *ReturnStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitReturnStmt(s)
}

type ClassStmt struct {
	Name    token.Token
	Methods []Stmt
}

func (s *ClassStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitClassStmt(s)
}
