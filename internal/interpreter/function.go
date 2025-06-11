package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
)

type LoxFunction struct {
	declaration   ast.FunctionStmt
	closure       Environment
	isInitializer bool
}

func newLoxFunction(declaration ast.FunctionStmt, closure Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

type LoxFunctionReturnValue struct {
	Value any
}

func (lf *LoxFunction) Call(interpreter Interpreter, arguments []any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			if returnVal, ok := r.(LoxFunctionReturnValue); ok {
				result = returnVal.Value
				err = nil
			} else {
				panic(r)
			}
		}
		if lf.isInitializer {
			result, err = lf.closure.getAt(0, "this")
		}
	}()

	environment := newEnvironment(&lf.closure)

	for i, param := range lf.declaration.Parameters {
		environment.define(param.Lexeme, arguments[i])
	}

	return nil, interpreter.executeBlock(lf.declaration.Body, environment)
}

func (lf *LoxFunction) bind(instance instance) *LoxFunction {
	env := newEnvironment(&lf.closure)
	env.define("this", instance)
	return newLoxFunction(lf.declaration, env, lf.isInitializer)
}

func (lf *LoxFunction) Arity() int {
	return len(lf.declaration.Parameters)
}

func (lf *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", lf.declaration.Name.Lexeme)
}
