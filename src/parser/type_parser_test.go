package parser

import (
	"bananascript/src/lexer"
	"bananascript/src/types"
	"gotest.tools/assert"
	"testing"
)

func TestTypeParser(t *testing.T) {

	assertType(t, "string", &types.StringType{})
	assertType(t, "bool", &types.BoolType{})
	assertType(t, "int", &types.IntType{})
	assertType(t, "null", &types.NullType{})
	assertType(t, "void", &types.VoidType{})
	assertType(t, "doesNotExist", &types.NeverType{})

	assertType(t, "string?", &types.Optional{Base: &types.StringType{}})
	assertType(t, "int?", &types.Optional{Base: &types.IntType{}})
	assertType(t, "bool????", &types.Optional{Base: &types.BoolType{}})

	assertType(t,
		"fn(string, fn() void, bool?) int?",
		&types.FunctionType{
			ParameterTypes: []types.Type{
				&types.StringType{},
				&types.FunctionType{
					ParameterTypes: []types.Type{},
					ReturnType:     &types.VoidType{},
				},
				&types.Optional{Base: &types.BoolType{}},
			},
			ReturnType: &types.Optional{Base: &types.IntType{}},
		},
	)

	assertType(t,
		"fn string",
		&types.NeverType{},
	)
}

func assertType(t *testing.T, input string, expected types.Type) {

	theLexer := lexer.FromCode(input)
	parser := New(theLexer)

	context := NewContext()
	theType := parser.parseType(context, TypeLowest)

	assert.DeepEqual(t, theType, expected)
}
