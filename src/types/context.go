package types

import (
	"reflect"
)

type Context struct {
	parent       *Context
	typeContexts map[Type]*Context
	memberStore  map[string]Type
	typeStore    map[string]Type
	ReturnType   Type
}

func NewContext() *Context {
	return &Context{
		typeContexts: make(map[Type]*Context),
		memberStore:  make(map[string]Type),
		typeStore:    make(map[string]Type),
	}
}

func ExtendContext(parent *Context) *Context {
	return &Context{
		parent:       parent,
		ReturnType:   parent.ReturnType,
		typeContexts: make(map[Type]*Context),
		memberStore:  make(map[string]Type),
		typeStore:    make(map[string]Type),
	}
}

func GetMemberTypeContext(context *Context, parentType Type) *Context {
	newContext := NewContext()
	currentContext := context
	for currentContext != nil && parentType != nil {
		for memberParentType, typeContext := range currentContext.typeContexts {
			if memberParentType.IsAssignable(parentType, context) {
				for memberName, memberType := range typeContext.memberStore {
					newContext.memberStore[memberName] = memberType
				}
			}
		}
		currentContext = currentContext.parent
	}
	if iface, isIface := parentType.(*Iface); isIface {
		for memberName, memberType := range iface.Members {
			newContext.memberStore[memberName] = memberType
		}
	}
	return newContext
}

func CloneContext(context *Context) *Context {
	return &Context{
		parent:       context.parent,
		ReturnType:   context.ReturnType,
		typeContexts: cloneTypeMap(context.typeContexts),
		memberStore:  cloneMap(context.memberStore),
		typeStore:    cloneMap(context.typeStore),
	}
}

func (context *Context) GetMemberTypeStrict(name string) (Type, bool) {
	memberType, ok := context.memberStore[name]
	return memberType, ok
}

func (context *Context) GetMemberType(name string) (Type, bool) {
	memberType, ok := context.GetMemberTypeStrict(name)
	if !ok && context.parent != nil {
		return context.parent.GetMemberType(name)
	}
	return memberType, ok
}

func (context *Context) DefineMemberType(name string, memberType Type) (Type, bool) {
	if _, exists := context.GetMemberTypeStrict(name); exists {
		return nil, false
	}
	context.memberStore[name] = memberType
	return memberType, true
}

func (context *Context) GetTypeMemberTypeStrict(name string, parentType Type) (Type, Type, bool) {
	for resolvedParentType, typeContext := range context.typeContexts {
		memberType, ok := typeContext.GetMemberTypeStrict(name)
		if ok && resolvedParentType.IsAssignable(parentType, context) {
			return memberType, resolvedParentType, true
		}
	}
	return nil, nil, false
}

func (context *Context) GetTypeMemberType(name string, parentType Type) (Type, Type, bool) {
	memberType, resolvedParentType, ok := context.GetTypeMemberTypeStrict(name, parentType)
	if !ok && context.parent != nil {
		return context.parent.GetTypeMemberType(name, parentType)
	}
	return memberType, resolvedParentType, ok
}

func (context *Context) DefineTypeMemberType(name string, memberType Type, parentType Type) (Type, bool) {
	_, _, exists := context.GetTypeMemberTypeStrict(name, parentType)
	if exists {
		return nil, false
	}
	var parentTypeContext *Context
	for theParentType, typeContext := range context.typeContexts {
		if reflect.DeepEqual(parentType, theParentType) {
			parentTypeContext = typeContext
			break
		}
	}
	if parentTypeContext == nil {
		parentTypeContext = NewContext()
		context.typeContexts[parentType] = parentTypeContext
	}
	return parentTypeContext.DefineMemberType(name, memberType)
}

func (context *Context) GetTypeStrict(name string) (Type, bool) {
	theType, ok := context.typeStore[name]
	return theType, ok
}

func (context *Context) GetType(name string) (Type, bool) {
	theType, ok := context.GetTypeStrict(name)
	if !ok && context.parent != nil {
		return context.parent.GetType(name)
	}
	return theType, ok
}

func (context *Context) DefineType(name string, theType Type) (Type, bool) {
	if _, exists := context.GetTypeStrict(name); exists {
		return nil, false
	}
	context.typeStore[name] = theType
	return theType, true
}

// TODO wait for generic fix
func cloneTypeMap[V any](toClone map[Type]V) map[Type]V {
	cloned := make(map[Type]V)
	for key, value := range toClone {
		cloned[key] = value
	}
	return cloned
}

func cloneMap[K comparable, V any](toClone map[K]V) map[K]V {
	cloned := make(map[K]V)
	for key, value := range toClone {
		cloned[key] = value
	}
	return cloned
}
