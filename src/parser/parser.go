package parser

import (
	"bananascript/src/errors"
	"bananascript/src/lexer"
	"bananascript/src/token"
	"bananascript/src/types"
	"reflect"
)

type Parser struct {
	errors   []*errors.ParserError
	tokens   []*token.Token
	position int
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

	parser := &Parser{tokens: tokens, errors: lexer.Errors}
	parser.registerExpressionParseFunctions()
	parser.registerTypeParseFunctions()
	return parser
}

func (parser *Parser) error(token *token.Token, messageFormat string, args ...interface{}) {
	parser.errors = append(parser.errors, errors.NewFromToken(token, messageFormat, args...))
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

func (parser *Parser) ParseProgram(context *types.Context) (*Program, []*errors.ParserError) {

	program := &Program{}
	program.Statements = []Statement{}
	program.Context = types.ExtendContext(context)

	for parser.current().Type != token.EOF {
		if parser.current().Type == token.Semi || parser.current().Type == token.Illegal {
			parser.consume()
			continue
		}

		statement := parser.parseStatement(program.Context)

		if statement != nil && !reflect.ValueOf(statement).IsNil() {
			program.Statements = append(program.Statements, statement)
		}
		parser.consume()
	}

	parser.doesReturn(context, program)
	return program, parser.errors
}

func (parser *Parser) doesReturn(context *types.Context, statement Statement) bool {

	switch statement := statement.(type) {
	case *Program:
		for _, statement := range statement.Statements {
			parser.doesReturn(context, statement)
		}
	case *ReturnStatement:
		if context.ReturnType != nil {
			if !isNever(context.ReturnType) {
				returnType := parser.getExpressionType(statement.Expression, context)
				if !isNever(returnType) && !context.ReturnType.IsAssignable(returnType, context) {
					parser.error(statement.ReturnToken, "Type '%s' is not assignable to '%s'", returnType.ToString(),
						context.ReturnType.ToString())
				}
			}
			return true
		} else {
			parser.error(statement.ReturnToken, "Illegal return statement")
		}
	case *BlockStatement:
		newContext := statement.Context
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
		return parser.doesReturn(statement.StatementContext, statement.Statement) &&
			parser.doesReturn(statement.AlternativeContext, statement.Alternative)
	}
	return false
}

func (parser *Parser) parseParameterList(context *types.Context) []*Parameter {

	parameters := make([]*Parameter, 0)
	if parser.peek().Type == token.RParen {
		parser.consume()
		return parameters
	}

	for {
		parameter := parser.parseParameter(context)
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

func (parser *Parser) parseParameter(context *types.Context) *Parameter {

	if !parser.assertNext(token.Ident) {
		return nil
	}
	identToken := parser.current()
	ident := &Identifier{IdentToken: identToken, Value: identToken.Literal}

	if !parser.assertNext(token.Colon) {
		return nil
	}

	parser.consume()
	theType := parser.parseType(context, TypeLowest)

	return &Parameter{Token: identToken, Name: ident, Type: theType}
}
