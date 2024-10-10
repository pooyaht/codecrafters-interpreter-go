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
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v Visitor) (any, error) {
	return v.VisitLiteralExpr(e)
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
