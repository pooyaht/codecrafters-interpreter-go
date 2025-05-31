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
