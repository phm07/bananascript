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

var prefixTypeParseFunctions = make(map[token.Type]func(*types.Context) types.Type)
var infixTypeParseFunctions = make(map[token.Type]func(*types.Context, types.Type) types.Type)

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
	prefixTypeParseFunctions[token.Iface] = parser.parseIfaceTypeLiteral

	infixTypeParseFunctions[token.Qmark] = parser.parseOptionalTypeLiteral
}

func (parser *Parser) parseType(context *types.Context, precedence TypePrecedence) types.Type {
	currentToken := parser.current()
	prefixFunction := prefixTypeParseFunctions[currentToken.Type]
	if prefixFunction == nil {
		if currentToken.Type != token.Illegal {
			parser.error(currentToken, "Unexpected %s", currentToken.ToString())
		}
		return &types.Never{}
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

func (parser *Parser) parseTypeLiteral(context *types.Context) types.Type {
	currentToken := parser.current()

	switch currentToken.Type {
	case token.Ident:
		typeName := parser.current().Literal
		switch typeName {
		case types.TypeString:
			return &types.String{}
		case types.TypeBool:
			return &types.Bool{}
		case types.TypeInt:
			return &types.Int{}
		case types.TypeFloat:
			return &types.Float{}
		default:
			theType, ok := context.GetType(typeName)
			if !ok {
				parser.error(currentToken, "Unknown type '%s'", typeName)
				return &types.Never{}
			}
			return theType
		}
	case token.Null:
		return &types.Null{}
	case token.Void:
		return &types.Void{}
	}
	return &types.Never{}
}

func (parser *Parser) parseFunctionTypeLiteral(context *types.Context) types.Type {
	if !parser.assertNext(token.LParen) {
		return &types.Never{}
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
		return &types.Never{}
	}
	parser.consume()
	returnType := parser.parseType(context, TypeLowest)
	return &types.Function{
		ParameterTypes: parameterTypes,
		ReturnType:     returnType,
	}
}

func (parser *Parser) parseIfaceTypeLiteral(context *types.Context) types.Type {

	if !parser.assertNext(token.LBrace) {
		return &types.Never{}
	}

	iface := &types.Iface{Members: make(map[string]types.Type)}
	for parser.peek().Type == token.Ident {
		parser.consume()
		name := parser.current().Literal
		parser.assertNext(token.Colon)
		parser.consume()
		memberType := parser.parseType(context, TypeLowest)
		iface.Members[name] = memberType
		parser.assertNext(token.Semi)
	}

	parser.assertNext(token.RBrace)

	return iface
}

/** infix types **/

func (parser *Parser) parseOptionalTypeLiteral(_ *types.Context, left types.Type) types.Type {
	switch left.(type) {
	case *types.Never, *types.Null, *types.Void, *types.Optional:
		return left
	default:
		return &types.Optional{Base: left}
	}
}
