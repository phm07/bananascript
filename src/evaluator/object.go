package evaluator

import (
	"bananascript/src/parser"
	"bananascript/src/types"
	"strconv"
)

type ObjectType = string

type Object interface {
	ToString() string
	Type() types.Type
}

type ErrorObject struct {
	Message string
}

func (errorObject *ErrorObject) ToString() string {
	return "ERROR: " + errorObject.Message
}

func (*ErrorObject) Type() types.Type {
	return nil
}

type ReturnObject struct {
	Object Object
}

func (returnObject *ReturnObject) ToString() string {
	return returnObject.Object.ToString()
}

func (returnObject *ReturnObject) Type() types.Type {
	return returnObject.Object.Type()
}

type Function interface {
	Object
	Execute(arguments []Object) Object
	With(object Object) Function
}

type FunctionObject struct {
	Environment  *Environment
	Parameters   []*parser.Identifier
	Body         *parser.BlockStatement
	This         Object
	Context      *types.Context
	FunctionType types.Type
}

func (functionObject *FunctionObject) Execute(arguments []Object) Object {
	newEnvironment := ExtendEnvironment(functionObject.Environment, functionObject.Context)
	if functionObject.This != nil {
		newEnvironment.DefineObject("this", functionObject.This)
	}
	for i, argument := range arguments {
		name := functionObject.Parameters[i].Value
		_, ok := newEnvironment.DefineObject(name, argument)
		if !ok {
			return NewError("Parameter %s already exists", name)
		}
	}
	return Eval(functionObject.Body, newEnvironment)
}

func (functionObject *FunctionObject) Type() types.Type {
	return functionObject.FunctionType
}

func (functionObject *FunctionObject) With(object Object) Function {
	newFunction := *functionObject
	newFunction.This = object
	return &newFunction
}

func (*FunctionObject) ToString() string {
	return "[Function]"
}

type StringObject struct {
	Value string
}

func (stringObject *StringObject) ToString() string {
	return stringObject.Value
}

func (*StringObject) Type() types.Type {
	return &types.String{}
}

type IntegerObject struct {
	Value int64
}

func (integerObject *IntegerObject) ToString() string {
	return strconv.FormatInt(integerObject.Value, 10)
}

func (*IntegerObject) Type() types.Type {
	return &types.Int{}
}

type BooleanObject struct {
	Value bool
}

func (booleanObject *BooleanObject) ToString() string {
	return strconv.FormatBool(booleanObject.Value)
}

func (*BooleanObject) Type() types.Type {
	return &types.Bool{}
}

type NullObject struct {
}

func (*NullObject) ToString() string {
	return "null"
}

func (*NullObject) Type() types.Type {
	return &types.Null{}
}
