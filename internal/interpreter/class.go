package interpreter

type class struct {
	name       string
	superclass *class
	methods    map[string]*LoxFunction
}

func newClass(name string, superclass *class, methods map[string]*LoxFunction) class {
	return class{
		name,
		superclass,
		methods,
	}
}

func (cls class) String() string {
	return cls.name
}

func (cls class) Call(interpreter Interpreter, arguments []any) (result any, err error) {
	instance := newInstance(cls)
	if initMethod := cls.findMethod("init"); initMethod != nil {
		return initMethod.bind(instance).Call(interpreter, arguments)
	}
	return instance, nil
}

func (cls class) Arity() int {
	if initMethod := cls.findMethod("init"); initMethod != nil {
		return initMethod.Arity()
	}
	return 0
}

func (cls class) findMethod(name string) *LoxFunction {
	if method, ok := cls.methods[name]; ok {
		return method
	}
	if cls.superclass != nil {
		return cls.superclass.findMethod(name)
	}
	return nil
}
