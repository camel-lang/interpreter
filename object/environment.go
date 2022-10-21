package object

func NewEnclosedEnvironment(out *Environment) *Environment {
	env := NewEnvironment()
	env.outer = out
	return env
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok := e.outer.Get(name)
		return obj, ok
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
