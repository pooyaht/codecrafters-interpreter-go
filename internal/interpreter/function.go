package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
)

type LoxFunction struct {
	decleration ast.FunctionStmt
}

func newLoxFunction(decleration ast.FunctionStmt) *LoxFunction {
	return &LoxFunction{
		decleration,
	}
}

type LoxFunctionReturnValue struct {
	Value any
}

func (lf *LoxFunction) Call(interpreter *Interpreter, arguments []any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			if returnVal, ok := r.(LoxFunctionReturnValue); ok {
				result = returnVal.Value
				err = nil
				return
			}
			panic(r)
		}
	}()

	environment := newEnvironment(&interpreter.globals)
	for i, param := range lf.decleration.Parameters {
		environment.define(param.Lexeme, arguments[i])
	}

	return nil, interpreter.executeBlock(lf.decleration.Body, environment)
}

func (lf *LoxFunction) Arity() int {
	return len(lf.decleration.Parameters)
}

func (lf *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", lf.decleration.Name.Lexeme)
}
