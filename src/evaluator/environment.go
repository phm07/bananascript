package evaluator

import "bananascript/src/types"

type Environment struct {
	parent      *Environment
	store       map[string]Object
	typeMembers map[types.Type]map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), typeMembers: make(map[types.Type]map[string]Object)}
}

func ExtendEnvironment(parent *Environment) *Environment {
	return &Environment{parent: parent, store: make(map[string]Object), typeMembers: make(map[types.Type]map[string]Object)}
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

func (environment *Environment) GetMember(object Object, objectType types.Type, name string) (Object, bool) {
	for theType, typeStore := range environment.typeMembers {
		if theType.IsAssignable(objectType) {
			if member, exists := typeStore[name]; exists {
				return member, true
			}
		}
	}
	if environment.parent != nil {
		return environment.parent.GetMember(object, objectType, name)
	}
	return nil, false
}

func (environment *Environment) Define(name string, value Object) (Object, bool) {
	environment.store[name] = value
	return value, true
}

func (environment *Environment) DefineTypeMember(parentType types.Type, name string, member Object) (Object, bool) {
	if _, exists := environment.typeMembers[parentType]; !exists {
		environment.typeMembers[parentType] = make(map[string]Object)
	}
	environment.typeMembers[parentType][name] = member
	return member, true
}

func (environment *Environment) Assign(name string, value Object) (Object, bool) {
	if _, exists := environment.GetInThisScope(name); exists {
		return environment.Define(name, value)
	} else if environment.parent != nil {
		return environment.parent.Assign(name, value)
	}
	return nil, false
}
