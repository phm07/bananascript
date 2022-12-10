package builtins

import (
	"bananascript/src/evaluator"
	"bananascript/src/parser"
	"bananascript/src/types"
	"fmt"
	"strconv"
	"strings"
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
		"print": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{&types.StringType{}},
				ReturnType:     &types.VoidType{},
			},
			Object: &BuiltinFunction{
				Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
					fmt.Print(arguments[0].ToString())
					return nil
				},
			},
		},
		"prompt": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{&types.StringType{}},
				ReturnType:     &types.StringType{},
			},
			Object: &BuiltinFunction{
				Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
					fmt.Print(arguments[0].ToString())
					var input string
					_, err := fmt.Scanln(&input)
					if err != nil {
						panic(err)
					}
					return &evaluator.StringObject{Value: input}
				},
			},
		},
		"min": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{&types.IntType{}, &types.IntType{}},
				ReturnType:     &types.IntType{},
			},
			Object: &BuiltinFunction{
				Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
					a := arguments[0].(*evaluator.IntegerObject).Value
					b := arguments[1].(*evaluator.IntegerObject).Value
					min := a
					if b < a {
						min = b
					}
					return &evaluator.IntegerObject{Value: min}
				},
			},
		},
		"max": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{&types.IntType{}, &types.IntType{}},
				ReturnType:     &types.IntType{},
			},
			Object: &BuiltinFunction{
				Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
					a := arguments[0].(*evaluator.IntegerObject).Value
					b := arguments[1].(*evaluator.IntegerObject).Value
					max := a
					if a < b {
						max = b
					}
					return &evaluator.IntegerObject{Value: max}
				},
			},
		},
	},
	&types.IntType{}: {
		"toString": toStringBuiltin,
		"abs": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.IntType{},
			},
			Object: &BuiltinFunction{
				Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
					value := this.(*evaluator.IntegerObject).Value
					if value < 0 {
						value *= -1
					}
					return &evaluator.IntegerObject{Value: value}
				},
			},
		},
	},
	&types.BoolType{}: {
		"toString": toStringBuiltin,
	},
	&types.StringType{}: {
		"toString": toStringBuiltin,
		"length": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.IntType{},
			},
			Object: &BuiltinFunction{
				Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
					return &evaluator.IntegerObject{Value: int64(len([]rune(this.ToString())))}
				},
			},
		},
		"uppercase": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.StringType{},
			},
			Object: &BuiltinFunction{
				Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
					return &evaluator.StringObject{Value: strings.ToUpper(this.ToString())}
				},
			},
		},
		"lowercase": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.StringType{},
			},
			Object: &BuiltinFunction{
				Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
					return &evaluator.StringObject{Value: strings.ToLower(this.ToString())}
				},
			},
		},
		"parseInt": {
			Type: &types.FunctionType{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.IntType{},
			},
			Object: &BuiltinFunction{
				Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
					// TODO throw error if invalid
					value, _ := strconv.ParseInt(this.ToString(), 10, 64)
					return &evaluator.IntegerObject{Value: value}
				},
			},
		},
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
