package parser

import (
	"bananascript/src/lexer"
	"bananascript/src/token"
	"bananascript/src/types"
	"reflect"
)

type Parser struct {
	errors               []*Error
	tokens               []*token.Token
	position             int
	prefixParseFunctions map[token.Type]func(*Context) Expression
	infixParseFunctions  map[token.Type]func(*Context, Expression) Expression
}

func New(lexer *lexer.Lexer) *Parser {
	var tokens []*token.Token
	for {
		nextToken := lexer.NextToken()
		tokens = append(tokens, nextToken)
		if nextToken.Type == token.EOF {
			break
		}
	}

	parser := &Parser{tokens: tokens}
	parser.registerExpressionParseFunctions()
	return parser
}

func (parser *Parser) error(token *token.Token, messageFormat string, args ...interface{}) {
	parser.errors = append(parser.errors, NewError(token, messageFormat, args...))
}

func (parser *Parser) current() *token.Token {
	if parser.position < len(parser.tokens) {
		return parser.tokens[parser.position]
	} else {
		return parser.tokens[parser.position-1]
	}
}

func (parser *Parser) consume() *token.Token {
	currentToken := parser.current()
	parser.position++
	return currentToken
}

func (parser *Parser) peek() *token.Token {
	if parser.position+1 < len(parser.tokens) {
		return parser.tokens[parser.position+1]
	} else {
		return parser.tokens[len(parser.tokens)-1]
	}
}

func (parser *Parser) assertNext(tokenType token.Type) bool {
	if nextToken := parser.peek(); nextToken.Type == tokenType {
		parser.consume()
		return true
	} else {
		parser.error(nextToken, "Expected %s, got %s instead", tokenType.ToStringHumanReadable(),
			nextToken.ToString())
		return false
	}
}

func (parser *Parser) ParseProgram(context *Context) (*Program, []*Error) {

	program := &Program{}
	program.Statements = []Statement{}

	for parser.current().Type != token.EOF {
		if parser.current().Type == token.Semi {
			parser.consume()
			continue
		}

		statement := parser.parseStatement(context)

		if statement != nil && !reflect.ValueOf(statement).IsNil() {
			program.Statements = append(program.Statements, statement)
		}
		parser.consume()
	}

	parser.doesReturn(context, program)
	parser.errors = removeDuplicateErrors(parser.errors)
	return program, parser.errors
}

func (parser *Parser) doesReturn(context *Context, statement Statement) bool {

	switch statement := statement.(type) {
	case *Program:
		for _, statement := range statement.Statements {
			parser.doesReturn(context, statement)
		}
	case *ReturnStatement:
		if context.returnType != nil {
			if !isNever(context.returnType) {
				returnType := statement.Expression.Type(context)
				if !isNever(returnType) && !context.returnType.IsAssignable(returnType) {
					parser.error(statement.ReturnToken, "Type '%s' is not assignable to '%s'", returnType.ToString(),
						context.returnType.ToString())
				}
			}
			return true
		} else {
			parser.error(statement.ReturnToken, "Illegal return statement")
		}
	case *BlockStatement:
		newContext := ExtendContext(context)
		returned := false
		for _, statement := range statement.Statements {
			if returned {
				parser.error(statement.Token(), "Unreachable code")
				return true
			}
			returned = parser.doesReturn(newContext, statement)
		}
		return returned
	case *IfStatement:
		return parser.doesReturn(ExtendContext(context), statement.Statement) &&
			parser.doesReturn(ExtendContext(context), statement.Alternative)
	}
	return false
}

func (parser *Parser) parseParameterList() []*Parameter {

	parameters := make([]*Parameter, 0)
	if parser.peek().Type == token.RParen {
		parser.consume()
		return parameters
	}

	for {
		parameter := parser.parseParameter()
		if parameter == nil {
			return nil
		}
		parameters = append(parameters, parameter)
		if parser.peek().Type == token.Comma {
			parser.consume()
		} else {
			break
		}
	}

	if !parser.assertNext(token.RParen) {
		return nil
	}
	return parameters
}

func (parser *Parser) parseParameter() *Parameter {

	if !parser.assertNext(token.Ident) {
		return nil
	}
	identToken := parser.current()
	ident := &Identifier{IdentToken: identToken, Value: identToken.Literal}

	if !parser.assertNext(token.Colon) {
		return nil
	}

	theType := parser.parseType()
	if theType == nil {
		return nil
	}

	return &Parameter{Token: identToken, Name: ident, Type: theType}
}

func (parser *Parser) parseType() types.Type {

	next := parser.peek()
	var currentType types.Type

	switch next.Type {
	case token.Ident:
		parser.consume()
		typeName := parser.current().Literal
		switch typeName {
		case types.TypeString:
			currentType = &types.StringType{}
		case types.TypeBool:
			currentType = &types.BoolType{}
		case types.TypeInt:
			currentType = &types.IntType{}
		default:
			currentType = newNever("Unknown type '%s'", typeName)
		}
	case token.Null:
		parser.consume()
		currentType = &types.NullType{}
	case token.Void:
		parser.consume()
		currentType = &types.VoidType{}
	case token.Func:
		parser.consume()
		if !parser.assertNext(token.LParen) {
			return nil
		}
		parameterTypes := make([]types.Type, 0)
		if parser.peek().Type != token.RParen {
			for {
				parameterType := parser.parseType()
				if parameterType == nil {
					return nil
				}
				parameterTypes = append(parameterTypes, parameterType)
				if parser.peek().Type == token.Comma {
					parser.consume()
				} else {
					break
				}
			}
		}
		if !parser.assertNext(token.RParen) {
			return nil
		}
		returnType := parser.parseType()
		if returnType == nil {
			returnType = &types.VoidType{}
		}
		return &types.FunctionType{
			ParameterTypes: parameterTypes,
			ReturnType:     returnType,
		}
	default:
		return nil
	}

	if parser.peek().Type == token.Qmark {
		parser.consume()
		currentType = types.NewOptional(currentType)
	}

	return currentType
}
