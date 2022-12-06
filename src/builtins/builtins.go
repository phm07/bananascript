package builtins

import (
	"bananascript/src/evaluator"
	"bananascript/src/parser"
	"fmt"
)

type Builtin struct {
	Type   parser.Type
	Object evaluator.Object
}

var builtins = map[string]*Builtin{
	"println": {
		Type: &parser.FunctionType{
			ParameterTypes: []parser.Type{&parser.StringType{}},
		},
		Object: &evaluator.FunctionObject{
			Parameters: []*parser.Identifier{{Value: "str"}},
			Execute:    builtinPrintln,
		},
	},
}

func NewContextAndEnvironment() (*parser.Context, *evaluator.Environment) {
	context := parser.NewContext()
	environment := evaluator.NewEnvironment()
	for name, builtin := range builtins {
		context.Define(name, builtin.Type)
		environment.Define(name, builtin.Object)
	}
	return parser.ExtendContext(context), evaluator.ExtendEnvironment(environment)
}

func builtinPrintln(arguments []evaluator.Object) evaluator.Object {
	fmt.Println(arguments[0].ToString())
	return nil
}
