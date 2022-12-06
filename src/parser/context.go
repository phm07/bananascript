package parser

type Context struct {
	parent     *Context
	store      map[string]Type
	returnType Type
}

func NewContext() *Context {
	return &Context{store: make(map[string]Type)}
}

func ExtendContext(parent *Context) *Context {
	return &Context{parent: parent, store: make(map[string]Type), returnType: parent.returnType}
}

func CloneContext(context *Context) *Context {
	cloned := &Context{
		parent:     context.parent,
		returnType: context.returnType,
		store:      make(map[string]Type),
	}
	for name, theType := range context.store {
		cloned.store[name] = theType
	}
	return cloned
}

func (context *Context) GetInThisScope(name string) (Type, bool) {
	object, ok := context.store[name]
	return object, ok
}

func (context *Context) GetType(name string) (Type, bool) {
	theType, ok := context.GetInThisScope(name)
	if !ok && context.parent != nil {
		return context.parent.GetType(name)
	}
	return theType, ok
}

func (context *Context) Define(name string, theType Type) Type {
	context.store[name] = theType
	return theType
}
