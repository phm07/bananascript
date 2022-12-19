package parser

import (
	"bananascript/src/lexer"
	"bananascript/src/token"
	"bananascript/src/types"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"testing"
)

func TestExpressionParser(t *testing.T) {

	assertExpression(t,
		"5 + 5",
		&InfixExpression{
			Left:     &IntegerLiteral{Value: 5},
			Operator: token.Plus,
			Right:    &IntegerLiteral{Value: 5},
		},
	)

	assertExpression(t,
		"1 + 2 * 3",
		&InfixExpression{
			Left:     &IntegerLiteral{Value: 1},
			Operator: token.Plus,
			Right: &InfixExpression{
				Left:     &IntegerLiteral{Value: 2},
				Operator: token.Star,
				Right:    &IntegerLiteral{Value: 3},
			},
		},
	)

	assertExpression(t,
		"func(++x, !(a || b))",
		&CallExpression{
			Function: &Identifier{Value: "func"},
			Arguments: []Expression{
				&IncrementExpression{
					Operator: token.Increment,
					Name:     &Identifier{Value: "x"},
					Pre:      true,
				},
				&PrefixExpression{
					Operator: token.Bang,
					Expression: &InfixExpression{
						Left:     &Identifier{Value: "a"},
						Operator: token.LogicalOr,
						Right:    &Identifier{Value: "b"},
					},
				},
			},
		},
	)

	assertExpression(t,
		"myString = \"hello\"",
		&AssignmentExpression{
			Name:       &Identifier{Value: "myString"},
			Expression: &StringLiteral{Value: "hello"},
		},
	)

	assertExpression(t,
		"+2",
		&InvalidExpression{},
	)

	assertExpression(t,
		"$",
		&InvalidExpression{},
	)
}

func assertExpression(t *testing.T, input string, expected Expression) {

	theLexer := lexer.FromCode(input)
	parser := New(theLexer)

	context := types.NewContext()
	expression := parser.parseExpression(context, ExpressionLowest)

	ignoreTokens := cmp.Comparer(func(t1, t2 *token.Token) bool {
		return true
	})

	assert.DeepEqual(t, expression, expected, ignoreTokens)
}
