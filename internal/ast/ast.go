package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/util"
)

type Visitor interface {
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitBinaryExpr(*BinaryExpr) (any, error)
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

func (p *AstPrinter) VisitUnaryExpr(e *UnaryExpr) (any, error) {
	str, _ := e.Right.Accept(p)
	return fmt.Sprintf("(%v %v)", e.Operator.Lexeme, str), nil
}

func (p *AstPrinter) VisitBinaryExpr(e *BinaryExpr) (any, error) {
	leftStr, _ := e.Left.Accept(p)
	rightStr, _ := e.Right.Accept(p)
	return fmt.Sprintf("(%v %v %v)", e.Operator.Lexeme, leftStr, rightStr), nil
}

type EvaluateVisitor struct{}

func (p *EvaluateVisitor) VisitLiteralExpr(e *LiteralExpr) (any, error) {
	return e.Value, nil
}

func (p *EvaluateVisitor) VisitGroupingExpr(e *GroupingExpr) (any, error) {
	return e.Expr.Accept(p)
}

func (p *EvaluateVisitor) VisitUnaryExpr(e *UnaryExpr) (any, error) {
	return nil, nil
}

func (p *EvaluateVisitor) VisitBinaryExpr(e *BinaryExpr) (any, error) {
	return nil, nil
}
