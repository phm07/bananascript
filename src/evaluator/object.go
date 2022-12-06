package evaluator

import (
	"bananascript/src/parser"
	"strconv"
)

type ObjectType = string

type Object interface {
	ToString() string
}

type ErrorObject struct {
	Message string
}

func (errorObject *ErrorObject) ToString() string {
	return "ERROR: " + errorObject.Message
}

type ReturnObject struct {
	Object Object
}

func (returnObject *ReturnObject) ToString() string {
	return returnObject.Object.ToString()
}

type FunctionObject struct {
	Parameters []*parser.Identifier
	Execute    func(arguments []Object) Object
}

func (functionObject *FunctionObject) ToString() string {
	return "[Function]"
}

type StringObject struct {
	Value string
}

func (stringObject *StringObject) ToString() string {
	return stringObject.Value
}

type IntegerObject struct {
	Value int64
}

func (integerObject *IntegerObject) ToString() string {
	return strconv.FormatInt(integerObject.Value, 10)
}

type BooleanObject struct {
	Value bool
}

func (booleanObject *BooleanObject) ToString() string {
	return strconv.FormatBool(booleanObject.Value)
}
