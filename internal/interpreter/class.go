package interpreter

type class struct {
	name    string
	methods map[string]LoxFunction
}

func newClass(name string, methods map[string]LoxFunction) class {
	return class{
		name,
		methods,
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

func (cls class) findMethod(name string) *LoxFunction {
	if method, ok := cls.methods[name]; ok {
		return &method
	}
	return nil
}
