package types

const (
	TypeNever  = "never"
	TypeNull   = "null"
	TypeVoid   = "void"
	TypeString = "string"
	TypeInt    = "int"
	TypeBool   = "bool"
)

type Type interface {
	ToString() string
	IsAssignable(Type) bool
}

type NeverType struct {
	Message string
}

func (neverType *NeverType) ToString() string {
	return TypeNever
}

func (neverType *NeverType) IsAssignable(_ Type) bool {
	return false
}

type NullType struct {
}

func (nullType *NullType) ToString() string {
	return TypeNull
}

func (nullType *NullType) IsAssignable(other Type) bool {
	return other.ToString() == TypeNull
}

type VoidType struct {
}

func (voidType *VoidType) ToString() string {
	return TypeVoid
}

func (voidType *VoidType) IsAssignable(other Type) bool {
	return other.ToString() == TypeVoid
}

type IntType struct {
}

func (integerType *IntType) ToString() string {
	return TypeInt
}

func (integerType *IntType) IsAssignable(other Type) bool {
	return other.ToString() == TypeInt
}

type BoolType struct {
}

func (boolType *BoolType) ToString() string {
	return TypeBool
}

func (boolType *BoolType) IsAssignable(other Type) bool {
	return other.ToString() == TypeBool
}

type StringType struct {
}

func (stringType *StringType) ToString() string {
	return TypeString
}

func (stringType *StringType) IsAssignable(other Type) bool {
	return other.ToString() == TypeString
}

type FunctionType struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (functionType *FunctionType) ToString() string {
	result := "fn("
	for i, parameter := range functionType.ParameterTypes {
		if i > 0 {
			result += ", "
		}
		result += parameter.ToString()
	}
	return result + ") " + functionType.ReturnType.ToString()
}

func (functionType *FunctionType) IsAssignable(other Type) bool {
	if other, isFunction := other.(*FunctionType); isFunction {
		if len(functionType.ParameterTypes) == len(other.ParameterTypes) {
			for i := range functionType.ParameterTypes {
				if !functionType.ParameterTypes[i].IsAssignable(other.ParameterTypes[i]) {
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

func (optional *Optional) IsAssignable(other Type) bool {
	switch other.ToString() {
	case optional.Base.ToString(), optional.Base.ToString() + "?", TypeNull:
		return true
	default:
		return false
	}
}

func NewOptional(base Type) Type {
	switch base.(type) {
	case *NeverType, *NullType, *VoidType:
		return base
	default:
		return &Optional{Base: base}
	}
}

func IsAssignable(theType Type, parentType Type) bool {
	if theType == nil {
		if parentType != nil {
			return false
		}
	} else {
		if parentType == nil || !theType.IsAssignable(parentType) {
			return false
		}
	}
	return true
}
