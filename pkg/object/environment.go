package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(id string) (Object, bool) {
	obj, ok := e.store[id]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(id)
	}
	return obj, ok
}

func (e *Environment) Set(id string, obj Object) Object {
	e.store[id] = obj
	return obj
}
