package lexer

import (
	"bananascript/src/token"
	"gotest.tools/assert"
	"testing"
)

func TestLexer(t *testing.T) {

	assertTypes(t,
		"let a := 5;",
		[]token.Type{token.Let, token.Ident, token.Define, token.IntLiteral, token.Semi},
	)

	assertTypes(t,
		"(5 + 7) / (3 * 8)",
		[]token.Type{token.LParen, token.IntLiteral, token.Plus, token.IntLiteral, token.RParen, token.Slash, token.LParen,
			token.IntLiteral, token.Star, token.IntLiteral, token.RParen},
	)

	assertTypes(t,
		"fn test(x: string) string { return \"test \" + x; }",
		[]token.Type{token.Func, token.Ident, token.LParen, token.Ident, token.Colon, token.Ident, token.RParen, token.Ident,
			token.LBrace, token.Return, token.StringLiteral, token.Plus, token.Ident, token.Semi, token.RBrace},
	)

	assertToken(t,
		"\"hello \\n \\\" \\t world\"",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "hello \n \" \t world",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)

	assertToken(t,
		"\"hello \\",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "hello \\",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)

	assertToken(t,
		"\"\\123\"",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "S",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)

	assertToken(t,
		"\"\\x1F60A\"",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "😊",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)

	assertToken(t,
		"\"\\U0001F408\"",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "🐈",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)

	assertToken(t,
		"\"你好世界\"",
		&token.Token{
			Type:    token.StringLiteral,
			Literal: "你好世界",
			Line:    1,
			Col:     1,
			File:    nil,
		},
	)
}

func assertTypes(t *testing.T, input string, expectedTypes []token.Type) {

	lexer := FromCode(input)
	lexedTypes := make([]token.Type, 0)

	for nextToken := lexer.NextToken(); nextToken.Type != token.EOF; nextToken = lexer.NextToken() {
		lexedTypes = append(lexedTypes, nextToken.Type)
	}

	if len(lexedTypes) > len(expectedTypes) {
		t.Errorf("Got more tokens than expected (%d > %d)", len(lexedTypes), len(expectedTypes))
		return
	} else if len(lexedTypes) < len(expectedTypes) {
		t.Errorf("Got less tokens than expected (%d < %d)", len(lexedTypes), len(expectedTypes))
		return
	}

	for i, lexedType := range lexedTypes {
		assert.Equal(t, lexedType, expectedTypes[i])
	}
}

func assertToken(t *testing.T, input string, expected *token.Token) {
	lexer := FromCode(input)
	theToken := lexer.NextToken()
	assert.DeepEqual(t, theToken, expected)
}
