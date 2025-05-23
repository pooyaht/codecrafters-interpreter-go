package ast

import "github.com/codecrafters-io/interpreter-starter-go/internal/token"

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

type VarStatement struct {
	Name        token.Token
	Initializer Expr
}

func (s *VarStatement) Accept(v StmtVisitor) (any, error) {
	return v.VisitVarStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitBlockStmt(s)
}
