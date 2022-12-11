package parser

import "bananascript/src/types"

type SubContext map[string]types.Type

type ContextStore map[types.Type]SubContext

type TypeStore map[string]types.Type

type Context struct {
	parent     *Context
	store      ContextStore
	typeStore  TypeStore
	returnType types.Type
}

func NewContext() *Context {
	return &Context{store: make(ContextStore), typeStore: make(TypeStore)}
}

func ExtendContext(parent *Context) *Context {
	return &Context{parent: parent, store: make(ContextStore), typeStore: make(TypeStore), returnType: parent.returnType}
}

func NewSubContext(context *Context, parentType types.Type) *Context {
	newContext := NewContext()
	newContext.store[nil] = make(SubContext)
	currentContext := context
	for currentContext != nil {
		for theType, typeStore := range currentContext.store {
			if theType != nil && theType.IsAssignable(parentType) {
				for name, memberType := range typeStore {
					if _, exists := newContext.store[nil][name]; !exists {
						newContext.store[nil][name] = memberType
					}
				}
			}
		}
		currentContext = currentContext.parent
	}
	return newContext
}

func CloneContext(context *Context) *Context {
	cloned := &Context{
		parent:     context.parent,
		returnType: context.returnType,
		store:      make(ContextStore),
	}
	for parentType, typeStore := range context.store {
		cloned.store[parentType] = make(SubContext)
		for name, memberType := range typeStore {
			cloned.store[parentType][name] = memberType
		}
	}
	return cloned
}

func (context *Context) GetInThisScope(name string, parentType types.Type) (types.Type, bool) {
	for theType, typeStore := range context.store {
		if types.IsAssignable(theType, parentType) {
			memberType, ok := typeStore[name]
			if ok {
				return memberType, true
			}
		}
	}
	return nil, false
}

func (context *Context) Get(name string, parentType types.Type) (types.Type, bool) {
	theType, ok := context.GetInThisScope(name, parentType)
	if !ok && context.parent != nil {
		return context.parent.Get(name, parentType)
	}
	return theType, ok
}

func (context *Context) Define(name string, theType types.Type, parentType types.Type) (types.Type, bool) {
	_, exists := context.GetInThisScope(name, parentType)
	if exists {
		return nil, false
	}
	if _, exists := context.store[parentType]; !exists {
		context.store[parentType] = make(map[string]types.Type)
	}
	context.store[parentType][name] = theType
	return theType, true
}

func (context *Context) GetTypeInThisScope(name string) (types.Type, bool) {
	theType, ok := context.typeStore[name]
	return theType, ok
}

func (context *Context) GetType(name string) (types.Type, bool) {
	theType, ok := context.GetTypeInThisScope(name)
	if !ok && context.parent != nil {
		return context.parent.GetType(name)
	}
	return theType, ok
}

func (context *Context) DefineType(name string, theType types.Type) (types.Type, bool) {
	if _, exists := context.GetTypeInThisScope(name); exists {
		return nil, false
	}
	context.typeStore[name] = theType
	return theType, true
}
