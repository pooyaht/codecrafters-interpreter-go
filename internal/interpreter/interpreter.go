package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
)

type Interpreter struct {
	astNodes []ast.Expr
}

func NewInterpreter(astNodes []ast.Expr) *Interpreter {
	return &Interpreter{astNodes: astNodes}
}

func (i *Interpreter) Interpret() {
	evaluator := &ast.EvaluateVisitor{}
	for _, node := range i.astNodes {
		value, _ := node.Accept(evaluator)
		if value == nil {
			fmt.Println("nil")
		} else {
			fmt.Println(value)
		}
	}
}
