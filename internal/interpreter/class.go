package interpreter

type class struct {
	name string
}

func newClass(name string) class {
	return class{
		name,
	}
}

func (cls class) String() string {
	return cls.name
}

func (cls class) Call(interpreter Interpreter, arguments []any) (result any, err error) {
	instance := newInstance(cls)
	return instance, nil
}

func (cls class) Arity() int {
	return 0
}
