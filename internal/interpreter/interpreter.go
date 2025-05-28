package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type RuntimeError struct {
	Message string
	Line    int
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s \n[line %d]\n", e.Message, e.Line)
}

type Interpreter struct {
	environment      Environment
	globals          Environment
	isInsideFunction bool
}

func NewInterpreter() Interpreter {
	globals := newEnvironment(nil)
	globals.define("clock", &ClockFunction{})
	return Interpreter{
		environment:      globals,
		globals:          globals,
		isInsideFunction: false,
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

func (i *Interpreter) VisitVarStmt(s *ast.VarStmt) (any, error) {
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

func (i *Interpreter) VisitBlockStmt(s *ast.BlockStmt) (any, error) {
	return nil, i.executeBlock(s.Statements, i.environment)
}

func (i *Interpreter) VisitIfStmt(s *ast.IfStmt) (any, error) {
	val, err := s.Condition.Accept(i)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(val) {
		return s.ThenBranch.Accept(i)
	} else if s.ElseBranch != nil {
		return s.ElseBranch.Accept(i)
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(s *ast.WhileStmt) (any, error) {
	for {
		condition, err := s.Condition.Accept(i)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(condition) {
			break
		}

		if _, err := s.Body.Accept(i); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(s *ast.FunctionStmt) (any, error) {
	function := newLoxFunction(*s)
	i.environment.define(s.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	if !i.isInsideFunction {
		return nil, RuntimeError{
			Message: "can't return from top-level code",
			Line:    stmt.Keyword.Line,
		}
	}

	var value any = nil
	if stmt.Value != nil {
		var err error
		value, err = stmt.Value.Accept(i)
		if err != nil {
			return nil, err
		}
	}
	panic(LoxFunctionReturnValue{Value: value})
}

func (i *Interpreter) VisitVariableExpr(e *ast.VariableExpr) (any, error) {
	return i.environment.get(e.Name)
}

func (i *Interpreter) VisitAssignmentExpr(e *ast.AssignmentExpr) (any, error) {
	value, err := e.Value.Accept(i)
	if err != nil {
		return nil, err
	}
	i.environment.assign(e.Name, value)
	return value, nil
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

func (i *Interpreter) VisitLogicalExpr(e *ast.LogicalExpr) (any, error) {
	left, err := e.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	if e.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return e.Right.Accept(i)
}

func (i *Interpreter) VisitCallExpr(e *ast.CallExpr) (any, error) {
	callee, err := e.Callee.Accept(i)
	if err != nil {
		return nil, err
	}

	var args []any
	for _, arg := range e.Arguments {
		arg, err := arg.Accept(i)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	callable, ok := callee.(LoxCallable)
	if !ok {
		return nil, fmt.Errorf("function is not callable: %v", callee)
	}
	if callable.Arity() != len(args) {
		return nil, RuntimeError{
			Message: fmt.Sprintf("expected %d arguments but got %d", callable.Arity(), len(args)),
			Line:    e.Paren.Line,
		}
	}

	return callable.Call(i, args)
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

func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return RuntimeError{Message: "Operand must be a number.", Line: operator.Line}
	}
	return nil
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)

	if !leftOk || !rightOk {
		return RuntimeError{Message: "Operands must be numbers.", Line: operator.Line}
	}

	return nil
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment Environment) error {
	previousEnv := i.environment

	defer func() {
		i.environment = previousEnv
	}()

	i.environment = newEnvironment(&environment)

	for _, stmt := range statements {
		_, err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}
