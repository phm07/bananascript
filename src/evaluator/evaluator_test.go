package evaluator

import (
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"bananascript/src/types"
	"gotest.tools/assert"
	"testing"
)

func TestEvaluator(t *testing.T) {

	assertObject(t,
		"5 + 5;",
		&IntegerObject{Value: 10},
	)

	assertObject(t,
		"\"a\" + \"b\";",
		&StringObject{Value: "ab"},
	)

	assertObject(t,
		"!!0;",
		&BooleanObject{Value: false},
	)

	assertObject(t,
		"1 + 2 * 3 - 4;",
		&IntegerObject{Value: 3},
	)
}

func assertObject(t *testing.T, input string, expected Object) {

	theLexer := lexer.FromCode(input)
	theParser := parser.New(theLexer)

	context := types.NewContext()
	environment := NewEnvironment(context)
	program, errors := theParser.ParseProgram(context)

	if len(errors) > 0 {
		for _, err := range errors {
			t.Error(err.Message)
		}
	} else {
		assert.DeepEqual(t, Eval(program.Statements[0], environment), expected)
	}
}
