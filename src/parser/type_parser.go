package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
)

type TypePrecedence int

const (
	TypeLowest TypePrecedence = iota
	TypeOptional
)

var typePrecedences = map[token.Type]TypePrecedence{
	token.Qmark: TypeOptional,
}

var prefixTypeParseFunctions = make(map[token.Type]func(*Context) types.Type)
var infixTypeParseFunctions = make(map[token.Type]func(*Context, types.Type) types.Type)

func getTypePrecedence(token *token.Token) TypePrecedence {
	if precedence, exists := typePrecedences[token.Type]; exists {
		return precedence
	}
	return TypeLowest
}

func (parser *Parser) registerTypeParseFunctions() {
	prefixTypeParseFunctions[token.Ident] = parser.parseTypeLiteral
	prefixTypeParseFunctions[token.Null] = parser.parseTypeLiteral
	prefixTypeParseFunctions[token.Void] = parser.parseTypeLiteral
	prefixTypeParseFunctions[token.Func] = parser.parseFunctionTypeLiteral

	infixTypeParseFunctions[token.Qmark] = parser.parseOptionalTypeLiteral
}

func (parser *Parser) parseType(context *Context, precedence TypePrecedence) types.Type {
	currentToken := parser.current()
	prefixFunction := prefixTypeParseFunctions[currentToken.Type]
	if prefixFunction == nil {
		if currentToken.Type != token.Illegal {
			parser.error(currentToken, "Unexpected %s", currentToken.ToString())
		}
		return &types.NeverType{}
	}

	theType := prefixFunction(context)

	for parser.peek().Type != token.Semi && precedence < getTypePrecedence(parser.peek()) {
		infixFunction := infixTypeParseFunctions[parser.peek().Type]
		if infixFunction == nil {
			break
		}

		parser.consume()
		theType = infixFunction(context, theType)
	}

	return theType
}

/** prefix types **/

func (parser *Parser) parseTypeLiteral(context *Context) types.Type {
	currentToken := parser.current()

	switch currentToken.Type {
	case token.Ident:
		typeName := parser.current().Literal
		switch typeName {
		case types.TypeString:
			return &types.StringType{}
		case types.TypeBool:
			return &types.BoolType{}
		case types.TypeInt:
			return &types.IntType{}
		default:
			theType, ok := context.GetType(typeName)
			if !ok {
				parser.error(currentToken, "Unknown type '%s'", typeName)
				return &types.NeverType{}
			}
			return theType
		}
	case token.Null:
		return &types.NullType{}
	case token.Void:
		return &types.VoidType{}
	}
	return &types.NeverType{}
}

func (parser *Parser) parseFunctionTypeLiteral(context *Context) types.Type {
	if !parser.assertNext(token.LParen) {
		return &types.NeverType{}
	}
	parameterTypes := make([]types.Type, 0)
	if parser.peek().Type != token.RParen {
		for {
			parser.consume()
			parameterType := parser.parseType(context, TypeLowest)
			parameterTypes = append(parameterTypes, parameterType)
			if parser.peek().Type == token.Comma {
				parser.consume()
			} else {
				break
			}
		}
	}
	if !parser.assertNext(token.RParen) {
		return &types.NeverType{}
	}
	parser.consume()
	returnType := parser.parseType(context, TypeLowest)
	return &types.FunctionType{
		ParameterTypes: parameterTypes,
		ReturnType:     returnType,
	}
}

/** infix types **/

func (parser *Parser) parseOptionalTypeLiteral(_ *Context, left types.Type) types.Type {
	switch left.(type) {
	case *types.NeverType, *types.NullType, *types.VoidType, *types.Optional:
		return left
	default:
		return &types.Optional{Base: left}
	}
}
