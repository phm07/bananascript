package parser

import (
	"bananascript/src/lexer"
	"bananascript/src/token"
	"bananascript/src/types"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"testing"
)

func TestStatementParser(t *testing.T) {

	assertStatement(t,
		"let a := 5;",
		&LetStatement{
			Name:  &Identifier{Value: "a"},
			Type:  &types.Int{},
			Value: &IntegerLiteral{Value: 5},
		},
	)

	assertStatement(t,
		"if a == 5 || a == 3 println(\"test\");",
		&IfStatement{
			Condition: &InfixExpression{
				Left: &InfixExpression{
					Left:     &Identifier{Value: "a"},
					Operator: token.EQ,
					Right:    &IntegerLiteral{Value: 5},
				},
				Operator: token.LogicalOr,
				Right: &InfixExpression{
					Left:     &Identifier{Value: "a"},
					Operator: token.EQ,
					Right:    &IntegerLiteral{Value: 3},
				},
			},
			Statement: &ExpressionStatement{
				Expression: &CallExpression{
					Function: &Identifier{Value: "println"},
					Arguments: []Expression{
						&StringLiteral{Value: "test"},
					},
				},
			},
		},
	)

	assertStatement(t,
		"type optionalString := string?;",
		&TypeDefinitionStatement{
			Name: &Identifier{Value: "optionalString"},
			Type: &types.Optional{Base: &types.String{}},
		},
	)

	assertError(t, "true / false;")
	assertError(t, "if true * false {}")
	assertError(t, "while \"a\" - 2 {}")
	assertError(t, "fn test(noType) {}")
	assertError(t, "fn noReturn() string {}")
	assertError(t, "{ type test := iface { abc: fn() void; }; let a: test = 2; }")

	assertNoError(t, "{ type str := string; let a: str = \"test\"; }")
	assertNoError(t, "{ type test := iface { }; let a: test = 0; let b: test = \"\"; let c: test = false; }")
}

func assertStatement(t *testing.T, input string, expected Statement) {

	theLexer := lexer.FromCode(input)
	parser := New(theLexer)

	context := types.NewContext()
	statement := parser.parseStatement(context)

	ignoreTokens := cmp.Comparer(func(t1, t2 *token.Token) bool {
		return true
	})
	ignoreContext := cmp.Comparer(func(c1, c2 *types.Context) bool {
		return true
	})

	assert.DeepEqual(t, statement, expected, ignoreTokens, ignoreContext)
}

func parse(input string) *Parser {
	theLexer := lexer.FromCode(input)
	theParser := New(theLexer)
	context := types.NewContext()
	theParser.parseStatement(context)
	return theParser
}

func assertError(t *testing.T, input string) {
	theParser := parse(input)
	assert.Assert(t, len(theParser.errors) > 0)
}

func assertNoError(t *testing.T, input string) {

	theParser := parse(input)

	errorMessages := make([]string, len(theParser.errors))
	for i, err := range theParser.errors {
		errorMessages[i] = err.Message
	}

	assert.Assert(t, len(theParser.errors) == 0, "\ninput: %s\nerrors: %v", input, errorMessages)
}
