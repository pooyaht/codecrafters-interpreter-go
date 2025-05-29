package interpreter

type LoxCallable interface {
	Call(interpreter Interpreter, arguments []any) (any, error)
	Arity() int
	String() string
}
