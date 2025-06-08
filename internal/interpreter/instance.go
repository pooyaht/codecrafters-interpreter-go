package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type instance struct {
	class  class
	fields map[string]any
}

func (inst instance) get(name token.Token) (any, error) {
	if val, ok := inst.fields[name.Lexeme]; ok {
		return val, nil
	}

	return nil, RuntimeError{Message: fmt.Sprintf("undefined property %s", name.Lexeme), Line: name.Line}
}

func newInstance(class class) instance {
	return instance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (inst instance) String() string {
	return inst.class.name + " instance"
}
