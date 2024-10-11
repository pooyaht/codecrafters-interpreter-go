package ast

type Expr interface {
	Accept(Visitor) (any, error)
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
