package evaluator

type Environment struct {
	parent *Environment
	store  map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func ExtendEnvironment(parent *Environment) *Environment {
	return &Environment{parent: parent, store: make(map[string]Object)}
}

func (environment *Environment) GetInThisScope(name string) (Object, bool) {
	object, ok := environment.store[name]
	return object, ok
}

func (environment *Environment) Get(name string) (Object, bool) {
	object, ok := environment.GetInThisScope(name)
	if !ok && environment.parent != nil {
		return environment.parent.Get(name)
	}
	return object, ok
}

func (environment *Environment) Define(name string, value Object) (Object, bool) {
	if _, exists := environment.GetInThisScope(name); exists {
		return nil, false
	}
	environment.store[name] = value
	return value, true
}

func (environment *Environment) Assign(name string, value Object) (Object, bool) {
	if _, exists := environment.GetInThisScope(name); exists {
		environment.store[name] = value
		return value, true
	} else if environment.parent != nil {
		return environment.parent.Assign(name, value)
	}
	return nil, false
}
