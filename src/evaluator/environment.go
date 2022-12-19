package evaluator

import (
	"bananascript/src/types"
	"reflect"
)

type Environment struct {
	context          *types.Context
	parent           *Environment
	store            map[string]Object
	typeEnvironments map[types.Type]*Environment
}

func NewEnvironment(context *types.Context) *Environment {
	return &Environment{context: context, store: make(map[string]Object), typeEnvironments: make(map[types.Type]*Environment)}
}

func ExtendEnvironment(parent *Environment, context *types.Context) *Environment {
	return &Environment{context: context, parent: parent, store: make(map[string]Object), typeEnvironments: make(map[types.Type]*Environment)}
}

func (environment *Environment) GetObjectStrict(name string) (Object, bool) {
	object, ok := environment.store[name]
	return object, ok
}

func (environment *Environment) GetObject(name string) (Object, bool) {
	object, ok := environment.GetObjectStrict(name)
	if !ok && environment.parent != nil {
		return environment.parent.GetObject(name)
	}
	return object, ok
}

func (environment *Environment) GetTypeMember(object Object, parentType types.Type, name string) (Object, bool) {
	for theType, typeStore := range environment.typeEnvironments {
		if theType.IsAssignable(parentType, environment.context) {
			object, ok := typeStore.GetObject(name)
			if ok {
				return object, ok
			}
		}
	}
	if environment.parent != nil {
		return environment.parent.GetTypeMember(object, parentType, name)
	}
	return nil, false
}

func (environment *Environment) DefineObject(name string, value Object) (Object, bool) {
	environment.store[name] = value
	return value, true
}

func (environment *Environment) DefineTypeMember(parentType types.Type, name string, member Object) (Object, bool) {
	for theType, typeStore := range environment.typeEnvironments {
		if reflect.DeepEqual(theType, parentType) {
			return typeStore.DefineObject(name, member)
		}
	}
	typeContext := types.GetMemberTypeContext(environment.context, parentType)
	typeEnvironment := NewEnvironment(typeContext)
	typeEnvironment.DefineObject(name, member)
	environment.typeEnvironments[parentType] = typeEnvironment
	return member, true
}

func (environment *Environment) AssignObject(name string, value Object) (Object, bool) {
	if _, exists := environment.GetObjectStrict(name); exists {
		return environment.DefineObject(name, value)
	} else if environment.parent != nil {
		return environment.parent.AssignObject(name, value)
	}
	return nil, false
}
