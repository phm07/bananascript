package builtins

import (
	"bananascript/src/evaluator"
	"bananascript/src/parser"
	"bananascript/src/types"
	"fmt"
)

type Builtin struct {
	Type   types.Type
	Object evaluator.Object
}

type BuiltinFunction struct {
	Executor func(evaluator.Object, []evaluator.Object) evaluator.Object
	This     evaluator.Object
}

func (builtinFunction *BuiltinFunction) Execute(arguments []evaluator.Object) evaluator.Object {
	return builtinFunction.Executor(builtinFunction.This, arguments)
}

func (builtinFunction *BuiltinFunction) With(object evaluator.Object) evaluator.Function {
	newFunction := *builtinFunction
	newFunction.This = object
	return &newFunction
}

func (builtinFunction *BuiltinFunction) ToString() string {
	return "[Function]"
}

var toStringBuiltin = &Builtin{
	Type: &types.FunctionType{
		ParameterTypes: []types.Type{},
		ReturnType:     &types.StringType{},
	},
	Object: &BuiltinFunction{
		Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
			return &evaluator.StringObject{Value: this.ToString()}
		},
	},
}

var builtinTypeMembers = map[types.Type]map[string]*Builtin{
	nil: {
		"println": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{&types.StringType{}},
				ReturnType:     &types.VoidType{},
			},
			Object: &BuiltinFunction{
				Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
					fmt.Println(arguments[0].ToString())
					return nil
				},
			},
		},
	},
	&types.IntType{}: {
		"toString": toStringBuiltin,
	},
	&types.BoolType{}: {
		"toString": toStringBuiltin,
	},
	&types.StringType{}: {
		"toString": toStringBuiltin,
	},
}

func NewContextAndEnvironment() (*parser.Context, *evaluator.Environment) {
	context := parser.NewContext()
	environment := evaluator.NewEnvironment()
	for parentType, builtins := range builtinTypeMembers {
		for name, builtin := range builtins {
			context.Define(name, builtin.Type, parentType)
			if parentType == nil {
				environment.Define(name, builtin.Object)
			} else {
				environment.DefineTypeMember(parentType, name, builtin.Object)
			}
		}
	}
	return parser.ExtendContext(context), evaluator.ExtendEnvironment(environment)
}
