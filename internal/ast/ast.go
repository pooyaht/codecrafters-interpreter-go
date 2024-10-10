package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/util"
)

type Expr interface {
	Accept(Visitor) (any, error)
}

type Visitor interface {
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v Visitor) (any, error) {
	return v.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) Accept(v Visitor) (any, error) {
	return v.VisitGroupingExpr(e)
}

type AstPrinter struct {
}

func (p *AstPrinter) VisitLiteralExpr(e *LiteralExpr) (any, error) {
	if e.Value == nil {
		return "nil", nil
	}

	if num, ok := e.Value.(float64); ok {
		return util.FormatFloat(num), nil
	}

	return fmt.Sprintf("%v", e.Value), nil
}

func (p *AstPrinter) VisitGroupingExpr(e *GroupingExpr) (any, error) {
	str, _ := e.Expr.Accept(p)
	return fmt.Sprintf("(group %v)", str), nil
}
