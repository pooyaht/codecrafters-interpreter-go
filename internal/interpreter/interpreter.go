package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Interpreter struct {
}

func (i *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(e *ast.LiteralExpr) (any, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(e *ast.GroupingExpr) (any, error) {
	return e.Expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(e *ast.UnaryExpr) (any, error) {
	rightEval, _ := e.Right.Accept(i)

	switch e.Operator.Type {
	case token.MINUS:
		if err := i.checkNumberOperand(rightEval); err != nil {
			return nil, err
		}
		return -rightEval.(float64), nil
	case token.BANG:
		return !i.isTruthy(rightEval), nil
	default:
		return nil, fmt.Errorf("unknown operator: %v", e.Operator.Lexeme)
	}
}

func (i *Interpreter) VisitBinaryExpr(e *ast.BinaryExpr) (any, error) {
	leftEval, _ := e.Left.Accept(i)
	rightEval, _ := e.Right.Accept(i)

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

func (i *Interpreter) isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) checkNumberOperand(operand any) error {
	if _, ok := operand.(float64); !ok {
		return fmt.Errorf("Operand must be a number.")
	}
	return nil
}
