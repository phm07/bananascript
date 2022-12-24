package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
	"reflect"
)

func (parser *Parser) parseStatement(context *types.Context) Statement {
	switch parser.current().Type {
	case token.Let:
		return parser.parseLetStatement(context)
	case token.Return:
		return parser.parseReturnStatement(context)
	case token.Func:
		return parser.parseFunctionDefinitionStatement(context)
	case token.LBrace:
		return parser.parseBlockStatement(context)
	case token.If:
		return parser.parseIfStatement(context)
	case token.While:
		return parser.parseWhileStatement(context)
	case token.TypeDef:
		return parser.parseTypeDefinitionStatement(context)
	default:
		return parser.parseExpressionStatement(context)
	}
}

func (parser *Parser) parseExpressionStatement(context *types.Context) *ExpressionStatement {
	statement := &ExpressionStatement{}
	statement.Expression = parser.parseExpression(context, ExpressionLowest)
	parser.getExpressionType(statement.Expression, context) // check for errors

	if !isInvalid(statement.Expression) {
		parser.assertNext(token.Semi)
	}

	return statement
}

func (parser *Parser) parseReturnStatement(context *types.Context) *ReturnStatement {
	statement := &ReturnStatement{ReturnToken: parser.consume()}

	if parser.current().Type == token.Semi {
		statement.Expression = &VoidLiteral{}
		return statement
	}

	statement.Expression = parser.parseExpression(context, ExpressionLowest)
	parser.assertNext(token.Semi)
	return statement
}

func (parser *Parser) parseBlockStatement(context *types.Context) *BlockStatement {

	newContext := types.ExtendContext(context)
	openingBrace := parser.consume()
	statements := make([]Statement, 0)

	for parser.current().Type != token.EOF && parser.current().Type != token.RBrace {
		if parser.current().Type == token.Semi || parser.current().Type == token.Illegal {
			parser.consume()
			continue
		}

		statement := parser.parseStatement(newContext)
		if statement != nil && !reflect.ValueOf(statement).IsNil() {
			statements = append(statements, statement)
		}
		parser.consume()
	}

	rBraceToken := parser.current()
	if rBraceToken.Type != token.RBrace {
		parser.error(openingBrace, "Unclosed block")
		rBraceToken = nil
	}

	return &BlockStatement{Statements: statements, LBraceToken: openingBrace, RBraceToken: rBraceToken, Context: newContext}
}

func (parser *Parser) parseLetStatement(context *types.Context) *LetStatement {
	statement := &LetStatement{LetToken: parser.current()}
	if !parser.assertNext(token.Ident) {
		return nil
	}

	identToken := parser.current()
	name := identToken.Literal
	statement.Name = &Identifier{IdentToken: identToken, Value: name}
	assignmentToken := token.Define

	if parser.peek().Type == token.Colon {
		parser.consume()
		parser.consume()
		statement.Type = parser.parseType(context, TypeLowest)
		assignmentToken = token.Assign
	}

	if parser.peek().Type != token.Semi {
		if !parser.assertNext(assignmentToken) {
			return nil
		}
		parser.consume()
		statement.Value = parser.parseExpression(context, ExpressionLowest)
	} else {
		statement.Value = &NullLiteral{}
	}

	parser.assertNext(token.Semi)

	inferredType := parser.getExpressionType(statement.Value, context)
	if statement.Type == nil {
		statement.Type = inferredType
	} else if !statement.Type.IsAssignable(inferredType, context) {
		erroneousToken := statement.Value.Token()
		if erroneousToken == nil {
			erroneousToken = parser.current()
		}
		if !isNever(statement.Type) && !isNever(inferredType) {
			parser.error(erroneousToken, "Type '%s' is not assignable to '%s'", inferredType.ToString(),
				statement.Type.ToString())
		}
	}

	_, ok := context.DefineMemberType(name, statement.Type)
	if !ok {
		parser.error(identToken, "Cannot redefine '%s'", name)
	}
	return statement
}

func (parser *Parser) parseFunctionDefinitionStatement(context *types.Context) *FunctionDefinitionStatement {

	statement := &FunctionDefinitionStatement{FuncToken: parser.current()}

	if parser.peek().Type == token.LParen {
		parser.consume() // fn
		parser.consume() // (
		statement.ThisType = parser.parseType(context, TypeLowest)
		if !parser.assertNext(token.RParen) || !parser.assertNext(token.DoubleColon) {
			return nil
		}
	}

	if !parser.assertNext(token.Ident) {
		return nil
	}
	identToken := parser.current()
	name := identToken.Literal
	statement.Name = &Identifier{IdentToken: identToken, Value: name}

	if !parser.assertNext(token.LParen) {
		return nil
	}

	statement.Parameters = parser.parseParameterList(context)
	if statement.Parameters == nil {
		return nil
	}
	parser.consume()

	if parser.current().Type == token.LBrace {
		statement.ReturnType = &types.Void{}
	} else {
		statement.ReturnType = parser.parseType(context, TypeLowest)
		if !parser.assertNext(token.LBrace) {
			return nil
		}
	}

	parameterTypes := make([]types.Type, 0)
	functionContext := types.ExtendContext(context)
	functionContext.ReturnType = statement.ReturnType
	if statement.ThisType != nil {
		functionContext.DefineMemberType("this", statement.ThisType)
	}
	for _, parameter := range statement.Parameters {
		parameterTypes = append(parameterTypes, parameter.Type)
		_, ok := functionContext.DefineMemberType(parameter.Name.Value, parameter.Type)
		if !ok {
			parser.error(parameter.Token, "Cannot redefine '%s'", parameter.Name.Value)
		}
	}

	statement.FunctionType = &types.Function{
		ParameterTypes: parameterTypes,
		ReturnType:     statement.ReturnType,
	}

	var ok bool
	if statement.ThisType != nil {
		_, ok = context.DefineTypeMemberType(name, statement.FunctionType, statement.ThisType)
	} else {
		_, ok = context.DefineMemberType(name, statement.FunctionType)
	}

	if !ok {
		parser.error(identToken, "Cannot redefine '%s'", name)
	}

	statement.FunctionContext = types.CloneContext(functionContext)
	statement.Body = parser.parseBlockStatement(statement.FunctionContext)
	if statement.Body == nil {
		return nil
	}

	if _, isVoid := statement.ReturnType.(*types.Void); !isVoid {
		if returns := parser.doesReturn(types.CloneContext(functionContext), statement.Body); !returns {
			erroneousToken := statement.Body.RBraceToken
			if erroneousToken == nil {
				erroneousToken = statement.Body.LBraceToken
			}
			parser.error(erroneousToken, "Missing return statement")
		}
	}

	return statement
}

func (parser *Parser) parseIfStatement(context *types.Context) *IfStatement {

	statement := &IfStatement{IfToken: parser.consume()}

	statement.Condition = parser.parseExpression(context, ExpressionLowest)
	parser.getExpressionType(statement.Condition, context) // check type
	parser.consume()

	statement.StatementContext = types.ExtendContext(context)
	statement.Statement = parser.parseStatement(statement.StatementContext)

	if parser.peek().Type == token.Else {
		parser.consume()
		parser.consume()
		statement.AlternativeContext = types.ExtendContext(context)
		statement.Alternative = parser.parseStatement(statement.AlternativeContext)
	}

	return statement
}

func (parser *Parser) parseWhileStatement(context *types.Context) *WhileStatement {

	statement := &WhileStatement{WhileToken: parser.consume()}

	statement.Condition = parser.parseExpression(context, ExpressionLowest)
	parser.getExpressionType(statement.Condition, context) // check type
	parser.consume()

	statement.StatementContext = types.ExtendContext(context)
	statement.Statement = parser.parseStatement(statement.StatementContext)

	return statement
}

func (parser *Parser) parseTypeDefinitionStatement(context *types.Context) *TypeDefinitionStatement {

	if !parser.assertNext(token.Ident) {
		return nil
	}
	identToken := parser.current()
	name := identToken.Literal
	ident := &Identifier{IdentToken: identToken, Value: name}

	if !parser.assertNext(token.Define) {
		return nil
	}

	parser.consume()
	statement := &TypeDefinitionStatement{IdentToken: identToken, Name: ident}
	statement.Type = parser.parseType(context, TypeLowest)

	switch name {
	case types.TypeNull, types.TypeVoid, types.TypeString, types.TypeInt, types.TypeFloat, types.TypeBool:
		parser.error(identToken, "Cannot re-declare primitive '%s'", name)
	default:
		if _, ok := context.DefineType(name, statement.Type); !ok {
			parser.error(identToken, "Cannot re-declare type '%s'", name)
		}
	}

	parser.assertNext(token.Semi)
	return statement
}
