package callable

type interpreterInterface any

type LoxCallable interface {
	Call(interpreter interpreterInterface, arguments []any) (any, error)
	Arity() int
	String() string
}
