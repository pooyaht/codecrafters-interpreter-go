package interpreter

type instance struct {
	class class
}

func newInstance(class class) instance {
	return instance{
		class: class,
	}
}

func (inst instance) String() string {
	return inst.class.name + " instance"
}
