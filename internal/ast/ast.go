package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/util"
)

type ExprVisitor interface {
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitVariableExpr(*VariableExpr) (any, error)
	VisitAssignmentExpr(*AssignmentExpr) (any, error)
}

type StmtVisitor interface {
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitExpressionStmt(*ExpressionStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitBlockStmt(*BlockStmt) (any, error)
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

func (p *AstPrinter) VisitVariableExpr(e *VariableExpr) (any, error) {
	return e.Name.Lexeme, nil
}

func (p *AstPrinter) VisitAssignmentExpr(e *AssignmentExpr) (any, error) {
	valueStr, _ := e.Value.Accept(p)
	return fmt.Sprintf("(= %v %v)", e.Name.Lexeme, valueStr), nil
}
