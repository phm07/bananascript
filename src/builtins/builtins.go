package builtins

import (
	"bananascript/src/evaluator"
	"bananascript/src/types"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BuiltinFunction struct {
	Executor     func(evaluator.Object, []evaluator.Object) evaluator.Object
	This         evaluator.Object
	FunctionType types.Type
}

func (builtinFunction *BuiltinFunction) Type() types.Type {
	return builtinFunction.FunctionType
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

var anyBuiltin = &types.Iface{
	Members: make(map[string]types.Type),
}

var builtinTypes = map[string]types.Type{
	"any": anyBuiltin,
}

var builtinObjects = map[types.Type]map[string]evaluator.Object{
	nil: {
		"println": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{anyBuiltin},
				ReturnType:     &types.Void{},
			},
			Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
				fmt.Println(arguments[0].ToString())
				return nil
			},
		},
		"print": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{anyBuiltin},
				ReturnType:     &types.Void{},
			},
			Executor: func(_ evaluator.Object, arguments []evaluator.Object) evaluator.Object {
				fmt.Print(arguments[0].ToString())
				return nil
			},
		},
		"prompt": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{anyBuiltin},
				ReturnType:     &types.String{},
			},
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
		"min": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{&types.Int{}, &types.Int{}},
				ReturnType:     &types.Int{},
			},
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
		"max": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{&types.Int{}, &types.Int{}},
				ReturnType:     &types.Int{},
			},
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
	anyBuiltin: {
		"toString": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.String{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.StringObject{Value: this.ToString()}
			},
		},
	},
	&types.Int{}: {
		"abs": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Int{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				value := this.(*evaluator.IntegerObject).Value
				if value < 0 {
					value *= -1
				}
				return &evaluator.IntegerObject{Value: value}
			},
		},
	},
	&types.Float{}: {
		"abs": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Float{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.FloatObject{Value: math.Abs(this.(*evaluator.FloatObject).Value)}
			},
		},
		"floor": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Float{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.FloatObject{Value: math.Floor(this.(*evaluator.FloatObject).Value)}
			},
		},
		"ceil": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Float{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.FloatObject{Value: math.Ceil(this.(*evaluator.FloatObject).Value)}
			},
		},
		"round": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Float{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.FloatObject{Value: math.Round(this.(*evaluator.FloatObject).Value)}
			},
		},
	},
	&types.String{}: {
		"length": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Int{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.IntegerObject{Value: int64(len([]rune(this.ToString())))}
			},
		},
		"uppercase": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.String{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.StringObject{Value: strings.ToUpper(this.ToString())}
			},
		},
		"lowercase": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.String{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				return &evaluator.StringObject{Value: strings.ToLower(this.ToString())}
			},
		},
		"parseInt": &BuiltinFunction{
			FunctionType: &types.Function{
				ParameterTypes: []types.Type{},
				ReturnType:     &types.Int{},
			},
			Executor: func(this evaluator.Object, _ []evaluator.Object) evaluator.Object {
				// TODO throw error if invalid
				value, _ := strconv.ParseInt(this.ToString(), 10, 64)
				return &evaluator.IntegerObject{Value: value}
			},
		},
	},
}

func NewContextAndEnvironment() (*types.Context, *evaluator.Environment) {
	context := types.NewContext()
	environment := evaluator.NewEnvironment(context)
	for parentType, builtins := range builtinObjects {
		for name, builtin := range builtins {
			if parentType == nil {
				context.DefineMemberType(name, builtin.Type())
				environment.DefineObject(name, builtin)
			} else {
				context.DefineTypeMemberType(name, builtin.Type(), parentType)
				environment.DefineTypeMember(parentType, name, builtin)
			}
		}
	}
	for builtinTypeName, builtinType := range builtinTypes {
		context.DefineType(builtinTypeName, builtinType)
	}
	return context, environment
}
