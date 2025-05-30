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
	if enclosing == nil {
		return Environment{
			enclosing: nil,
			values:    make(map[string]any, 0),
		}
	}

	enclosingCopy := *enclosing
	return Environment{
		enclosing: &enclosingCopy,
		values:    make(map[string]any, 0),
	}
}

func (e *Environment) get(name token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, fmt.Errorf("undefined variable %s", name.Lexeme)
}

func (e *Environment) assign(name token.Token, value any) (any, error) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return nil, fmt.Errorf("undefined variable %s", name.Lexeme)
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) getAt(distance int, name string) (any, error) {
	env := e.ancestor(distance)
	return env.values[name], nil
}

func (e *Environment) assignAt(distance int, name token.Token, value any) error {
	env := e.ancestor(distance)
	env.values[name.Lexeme] = value
	return nil
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for range distance {
		environment = environment.enclosing
	}
	return environment
}
