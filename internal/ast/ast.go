package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
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
	rightEval, _ := e.Right.Accept(p)

	isTruthy := func(v any) bool {
		if v == nil {
			return false
		}
		if b, ok := v.(bool); ok {
			return b
		}
		return true
	}

	switch e.Operator.Type {
	case token.MINUS:
		return -rightEval.(float64), nil
	case token.BANG:
		return !isTruthy(rightEval), nil
	default:
		return nil, fmt.Errorf("unknown operator: %v", e.Operator.Lexeme)
	}
}

func (p *EvaluateVisitor) VisitBinaryExpr(e *BinaryExpr) (any, error) {
	leftEval, _ := e.Left.Accept(p)
	rightEval, _ := e.Right.Accept(p)

	switch e.Operator.Type {
	case token.STAR:
		return leftEval.(float64) * rightEval.(float64), nil
	case token.SLASH:
		return leftEval.(float64) / rightEval.(float64), nil
	case token.PLUS:
		if leftStr, leftOk := leftEval.(string); leftOk {
			if rightStr, rightOk := rightEval.(string); rightOk {
				return leftStr + rightStr, nil
			}
		}
		return leftEval.(float64) + rightEval.(float64), nil
	case token.MINUS:
		return leftEval.(float64) - rightEval.(float64), nil
	case token.GREATER:
		return leftEval.(float64) > rightEval.(float64), nil
	case token.GREATER_EQUAL:
		return leftEval.(float64) >= rightEval.(float64), nil
	case token.LESS:
		return leftEval.(float64) < rightEval.(float64), nil
	case token.LESS_EQUAL:
		return leftEval.(float64) <= rightEval.(float64), nil
	case token.EQUAL_EQUAL:
		return leftEval == rightEval, nil
	case token.BANG_EQUAL:
		return leftEval != rightEval, nil
	default:
		return nil, fmt.Errorf("unknown operator: %v", e.Operator.Lexeme)
	}
}
