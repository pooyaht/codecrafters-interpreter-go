package ast

import (
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGroupingExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnaryExpr(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinaryExpr(e)
}

type VariableExpr struct {
	Name  token.Token
	Value Expr
}

func (e *VariableExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitVariableExpr(e)
}

type AssignmentExpr struct {
	Name  token.Token
	Value Expr
}

func (e *AssignmentExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitAssignmentExpr(e)
}

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *LogicalExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLogicalExpr(e)
}

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (e *CallExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitCallExpr(e)
}

type GetExpr struct {
	Object Expr
	Name   token.Token
}

func (e *GetExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGetExpr(e)
}
