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

type Function interface {
	Object
	Execute(arguments []Object) Object
	With(object Object) Function
}

type FunctionObject struct {
	Environment *Environment
	Parameters  []*parser.Identifier
	Body        *parser.BlockStatement
	This        Object
}

func (functionObject *FunctionObject) Execute(arguments []Object) Object {
	newEnvironment := ExtendEnvironment(functionObject.Environment)
	if functionObject.This != nil {
		newEnvironment.Define("this", functionObject.This)
	}
	for i, argument := range arguments {
		name := functionObject.Parameters[i].Value
		_, ok := newEnvironment.Define(name, argument)
		if !ok {
			return NewError("Parameter %s already exists", name)
		}
	}
	return Eval(functionObject.Body, newEnvironment)
}

func (functionObject *FunctionObject) With(object Object) Function {
	newFunction := *functionObject
	newFunction.This = object
	return &newFunction
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
