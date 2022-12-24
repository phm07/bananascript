package parser

import (
	"bananascript/src/lexer"
	"bananascript/src/types"
	"gotest.tools/assert"
	"testing"
)

func TestTypeParser(t *testing.T) {

	assertType(t, "string", &types.String{})
	assertType(t, "bool", &types.Bool{})
	assertType(t, "int", &types.Int{})
	assertType(t, "null", &types.Null{})
	assertType(t, "void", &types.Void{})
	assertType(t, "float", &types.Float{})
	assertType(t, "doesNotExist", &types.Never{})

	assertType(t, "string?", &types.Optional{Base: &types.String{}})
	assertType(t, "int?", &types.Optional{Base: &types.Int{}})
	assertType(t, "float?", &types.Optional{Base: &types.Float{}})
	assertType(t, "bool????", &types.Optional{Base: &types.Bool{}})

	assertType(t,
		"fn(string, fn() void, bool?) int?",
		&types.Function{
			ParameterTypes: []types.Type{
				&types.String{},
				&types.Function{
					ParameterTypes: []types.Type{},
					ReturnType:     &types.Void{},
				},
				&types.Optional{Base: &types.Bool{}},
			},
			ReturnType: &types.Optional{Base: &types.Int{}},
		},
	)

	assertType(t,
		"fn string",
		&types.Never{},
	)
}

func assertType(t *testing.T, input string, expected types.Type) {

	theLexer := lexer.FromCode(input)
	parser := New(theLexer)

	context := types.NewContext()
	theType := parser.parseType(context, TypeLowest)

	assert.DeepEqual(t, theType, expected)
}
