package lexer

import (
	"bananascript/src/token"
	"gotest.tools/assert"
	"testing"
)

func TestLexer(t *testing.T) {

	assertTokens(t,
		"let a := 5;",
		[]token.Type{token.Let, token.Ident, token.Define, token.IntLiteral, token.Semi},
	)

	assertTokens(t,
		"(5 + 7) / (3 * 8)",
		[]token.Type{token.LParen, token.IntLiteral, token.Plus, token.IntLiteral, token.RParen, token.Slash, token.LParen,
			token.IntLiteral, token.Star, token.IntLiteral, token.RParen},
	)

	assertTokens(t,
		"fn test(x: string) string { return \"test \" + x; }",
		[]token.Type{token.Func, token.Ident, token.LParen, token.Ident, token.Colon, token.Ident, token.RParen, token.Ident,
			token.LBrace, token.Return, token.StringLiteral, token.Plus, token.Ident, token.Semi, token.RBrace},
	)
}

func assertTokens(t *testing.T, input string, expectedTypes []token.Type) {

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
