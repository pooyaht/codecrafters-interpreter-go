package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	cerror "github.com/codecrafters-io/interpreter-starter-go/internal/error"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Interpreter struct {
	environment Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitPrintStmt(s *ast.PrintStmt) (any, error) {
	value, err := s.Expr.Accept(i)
	if err != nil {
		return nil, err
	}
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(e *ast.VariableExpr) (any, error) {
	return i.environment.get(e.Name)
}

func (i *Interpreter) VisitVarStmt(s *ast.VarStatement) (any, error) {
	var value any = nil
	var err error

	if s.Initializer != nil {
		value, err = s.Initializer.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	i.environment.define(s.Name.Lexeme, value)

	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(s *ast.ExpressionStmt) (any, error) {
	_, err := s.Expr.Accept(i)
	return nil, err
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
		if err := i.checkNumberOperand(e.Operator, rightEval); err != nil {
			return nil, err
		}
		return -rightEval.(float64), nil
	case token.BANG:
		return !i.isTruthy(rightEval), nil
	default:
		return nil, fmt.Errorf("unknown operator: %v at line %v", e.Operator.Lexeme, e.Operator.Line)
	}
}

func (i *Interpreter) VisitBinaryExpr(e *ast.BinaryExpr) (any, error) {
	leftEval, _ := e.Left.Accept(i)
	rightEval, _ := e.Right.Accept(i)

	switch e.Operator.Type {
	case token.STAR:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) * rightEval.(float64), nil
	case token.SLASH:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) / rightEval.(float64), nil
	case token.PLUS:
		if leftStr, leftOk := leftEval.(string); leftOk {
			if rightStr, rightOk := rightEval.(string); rightOk {
				return leftStr + rightStr, nil
			}
		}

		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) + rightEval.(float64), nil
	case token.MINUS:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) - rightEval.(float64), nil
	case token.GREATER:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) > rightEval.(float64), nil
	case token.GREATER_EQUAL:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) >= rightEval.(float64), nil
	case token.LESS:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
		return leftEval.(float64) < rightEval.(float64), nil
	case token.LESS_EQUAL:
		if err := i.checkNumberOperands(e.Operator, leftEval, rightEval); err != nil {
			return nil, err
		}
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

func (i *Interpreter) checkNumberOperand(operator *token.Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return cerror.RuntimeError{Message: "Operand must be a number.", Line: operator.Line}
	}
	return nil
}

func (i *Interpreter) checkNumberOperands(operator *token.Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)

	if !leftOk || !rightOk {
		return cerror.RuntimeError{Message: "Operands must be numbers.", Line: operator.Line}
	}

	return nil
}
