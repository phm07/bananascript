package types

const (
	TypeNever  = "never"
	TypeNull   = "null"
	TypeVoid   = "void"
	TypeString = "string"
	TypeInt    = "int"
	TypeFloat  = "float"
	TypeBool   = "bool"
)

type Type interface {
	ToString() string
	IsAssignable(Type, *Context) bool
}

type Never struct {
}

func (neverType *Never) ToString() string {
	return TypeNever
}

func (neverType *Never) IsAssignable(Type, *Context) bool {
	return false
}

type Null struct {
}

func (nullType *Null) ToString() string {
	return TypeNull
}

func (nullType *Null) IsAssignable(other Type, _ *Context) bool {
	_, isNull := other.(*Null)
	return isNull
}

type Void struct {
}

func (voidType *Void) ToString() string {
	return TypeVoid
}

func (voidType *Void) IsAssignable(other Type, _ *Context) bool {
	_, isVoid := other.(*Void)
	return isVoid
}

type Int struct {
}

func (integerType *Int) ToString() string {
	return TypeInt
}

func (integerType *Int) IsAssignable(other Type, _ *Context) bool {
	_, isInt := other.(*Int)
	return isInt
}

type Float struct {
}

func (floatType *Float) ToString() string {
	return TypeFloat
}

func (floatType *Float) IsAssignable(other Type, _ *Context) bool {
	_, isFloat := other.(*Float)
	return isFloat
}

type Bool struct {
}

func (boolType *Bool) ToString() string {
	return TypeBool
}

func (boolType *Bool) IsAssignable(other Type, _ *Context) bool {
	_, isBool := other.(*Bool)
	return isBool
}

type String struct {
}

func (stringType *String) ToString() string {
	return TypeString
}

func (stringType *String) IsAssignable(other Type, _ *Context) bool {
	_, isString := other.(*String)
	return isString
}

type Function struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (functionType *Function) ToString() string {
	result := "fn("
	for i, parameter := range functionType.ParameterTypes {
		if i > 0 {
			result += ", "
		}
		result += parameter.ToString()
	}
	return result + ") " + functionType.ReturnType.ToString()
}

func (functionType *Function) IsAssignable(other Type, context *Context) bool {
	if other, isFunction := other.(*Function); isFunction {
		if len(functionType.ParameterTypes) == len(other.ParameterTypes) {
			for i := range functionType.ParameterTypes {
				if !functionType.ParameterTypes[i].IsAssignable(other.ParameterTypes[i], context) {
					return false
				}
			}
			return true
		}
	}
	return false
}

type Optional struct {
	Base Type
}

func (optional *Optional) ToString() string {
	return optional.Base.ToString() + "?"
}

func (optional *Optional) IsAssignable(other Type, context *Context) bool {
	switch other := other.(type) {
	case *Null:
		return true
	case *Optional:
		return optional.Base.IsAssignable(other.Base, context)
	default:
		return optional.Base.IsAssignable(other, context)
	}
}

type Iface struct {
	Members map[string]Type
}

func (iface *Iface) ToString() string {
	result := "iface { "
	for name, memberType := range iface.Members {
		result += name + ": " + memberType.ToString() + "; "
	}
	return result + "}"
}

func (iface *Iface) IsAssignable(other Type, context *Context) bool {
	for name, memberType := range iface.Members {
		actualType, _, ok := context.GetTypeMemberType(name, other)
		if !ok || !memberType.IsAssignable(actualType, context) {
			return false
		}
	}
	return true
}
