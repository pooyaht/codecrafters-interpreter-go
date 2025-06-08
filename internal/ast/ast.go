package ast

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/internal/util"
)

type ExprVisitor interface {
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitVariableExpr(*VariableExpr) (any, error)
	VisitAssignmentExpr(*AssignmentExpr) (any, error)
	VisitLogicalExpr(*LogicalExpr) (any, error)
	VisitCallExpr(*CallExpr) (any, error)
	VisitGetExpr(*GetExpr) (any, error)
}

type StmtVisitor interface {
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitExpressionStmt(*ExpressionStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitBlockStmt(*BlockStmt) (any, error)
	VisitIfStmt(*IfStmt) (any, error)
	VisitWhileStmt(*WhileStmt) (any, error)
	VisitFunctionStmt(*FunctionStmt) (any, error)
	VisitReturnStmt(*ReturnStmt) (any, error)
	VisitClassStmt(*ClassStmt) (any, error)
}

type AstPrinter struct {
}

func (p *AstPrinter) VisitPrintStmt(s *PrintStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitExpressionStmt(s *ExpressionStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitVarStmt(s *VarStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitBlockStmt(s *BlockStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitIfStmt(s *BlockStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitWhileStmt(s *WhileStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitFunctionStmt(s *FunctionStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitClassStmt(s *ClassStmt) (any, error) {
	return nil, nil
}

func (p *AstPrinter) VisitLiteralExpr(e *LiteralExpr) (any, error) {
	if e.Value == nil {
		return "nil", nil
	}

	if num, ok := e.Value.(float64); ok {
		return util.FormatFloat(num, "parse"), nil
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

func (p *AstPrinter) VisitVariableExpr(e *VariableExpr) (any, error) {
	return e.Name.Lexeme, nil
}

func (p *AstPrinter) VisitAssignmentExpr(e *AssignmentExpr) (any, error) {
	valueStr, _ := e.Value.Accept(p)
	return fmt.Sprintf("(= %v %v)", e.Name.Lexeme, valueStr), nil
}

func (p *AstPrinter) VisitLogicalExpr(e *LogicalExpr) (any, error) {
	leftStr, _ := e.Left.Accept(p)
	rightStr, _ := e.Right.Accept(p)
	return fmt.Sprintf("(%v %v %v)", e.Operator.Lexeme, leftStr, rightStr), nil
}

func (p *AstPrinter) VisitCallExpr(e *CallExpr) (any, error) {
	callee, err := e.Callee.Accept(p)
	if err != nil {
		return nil, err
	}

	var args []string
	for _, arg := range e.Arguments {
		argStr, err := arg.Accept(p)
		if err != nil {
			return nil, err
		}
		args = append(args, argStr.(string))
	}

	return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", ")), nil
}

func (p *AstPrinter) VisitGetExpr(e *GetExpr) (any, error) {
	val, err := e.Object.Accept(p)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("(. %v %s)", val, e.Name.Lexeme), nil
}
