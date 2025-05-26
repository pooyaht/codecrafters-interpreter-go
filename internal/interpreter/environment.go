package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func newEnvironment(enclosing *Environment) Environment {
	return Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) get(name token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, fmt.Errorf("undefined varialbe %s", name.Lexeme)
}

func (e *Environment) assign(name token.Token, value any) (any, error) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil, nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return nil, fmt.Errorf("undefined varialbe %s", name.Lexeme)
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}
