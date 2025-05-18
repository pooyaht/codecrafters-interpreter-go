package ast

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type PrintStmt struct {
	Expr Expr
}

type ExpressionStmt struct {
	Expr Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrintStmt(s)
}
