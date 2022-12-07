package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
	"strconv"
)

type Precedence int

const (
	Lowest Precedence = iota
	Assignment
	LogicalOr
	LogicalAnd
	Equals
	Relation
	Sum
	Product
	Prefix
	Postfix
)

var precedences = map[token.Type]Precedence{
	token.Assign:     Assignment,
	token.LogicalOr:  LogicalOr,
	token.LogicalAnd: LogicalAnd,
	token.EQ:         Equals,
	token.NEQ:        Equals,
	token.LT:         Relation,
	token.GT:         Relation,
	token.LTE:        Relation,
	token.GTE:        Relation,
	token.Plus:       Sum,
	token.Minus:      Sum,
	token.Slash:      Product,
	token.Star:       Product,
	token.Increment:  Postfix,
	token.Decrement:  Postfix,
	token.LParen:     Postfix,
	token.Dot:        Postfix,
}

func getPrecedence(token *token.Token) Precedence {
	if precedence, exists := precedences[token.Type]; exists {
		return precedence
	}
	return Lowest
}

func (parser *Parser) registerExpressionParseFunctions() {
	parser.prefixParseFunctions = make(map[token.Type]func(*Context) Expression)
	parser.prefixParseFunctions[token.Ident] = parser.parseIdentifier
	parser.prefixParseFunctions[token.IntLiteral] = parser.parseIntegerLiteral
	parser.prefixParseFunctions[token.StringLiteral] = parser.parseStringLiteral
	parser.prefixParseFunctions[token.Null] = parser.parseNullLiteral
	parser.prefixParseFunctions[token.Void] = parser.parseVoidLiteral
	parser.prefixParseFunctions[token.True] = parser.parseBooleanLiteral
	parser.prefixParseFunctions[token.False] = parser.parseBooleanLiteral
	parser.prefixParseFunctions[token.Bang] = parser.parsePrefixExpression
	parser.prefixParseFunctions[token.Minus] = parser.parsePrefixExpression
	parser.prefixParseFunctions[token.LParen] = parser.parseGroupedExpression
	parser.prefixParseFunctions[token.Increment] = parser.parseIncrementPrefixExpression
	parser.prefixParseFunctions[token.Decrement] = parser.parseIncrementPrefixExpression

	parser.infixParseFunctions = make(map[token.Type]func(*Context, Expression) Expression)
	parser.infixParseFunctions[token.Assign] = parser.parseAssignmentExpression
	parser.infixParseFunctions[token.LogicalOr] = parser.parseInfixExpression
	parser.infixParseFunctions[token.LogicalAnd] = parser.parseInfixExpression
	parser.infixParseFunctions[token.EQ] = parser.parseInfixExpression
	parser.infixParseFunctions[token.NEQ] = parser.parseInfixExpression
	parser.infixParseFunctions[token.GT] = parser.parseInfixExpression
	parser.infixParseFunctions[token.LT] = parser.parseInfixExpression
	parser.infixParseFunctions[token.GTE] = parser.parseInfixExpression
	parser.infixParseFunctions[token.LTE] = parser.parseInfixExpression
	parser.infixParseFunctions[token.Plus] = parser.parseInfixExpression
	parser.infixParseFunctions[token.Minus] = parser.parseInfixExpression
	parser.infixParseFunctions[token.Slash] = parser.parseInfixExpression
	parser.infixParseFunctions[token.Star] = parser.parseInfixExpression
	parser.infixParseFunctions[token.LParen] = parser.parseCallExpression
	parser.infixParseFunctions[token.Increment] = parser.parseIncrementInfixExpression
	parser.infixParseFunctions[token.Decrement] = parser.parseIncrementInfixExpression
	parser.infixParseFunctions[token.Dot] = parser.parseMemberAccessExpression
}

func (parser *Parser) parseExpression(context *Context, precedence Precedence) Expression {
	currentToken := parser.current()
	prefixFunction := parser.prefixParseFunctions[currentToken.Type]
	if prefixFunction == nil {
		parser.error(currentToken, "Unexpected %s", currentToken.ToString())
		return &InvalidExpression{currentToken}
	}

	expression := prefixFunction(context)

	for parser.peek().Type != token.Semi && precedence < getPrecedence(parser.peek()) {
		infixFunction := parser.infixParseFunctions[parser.peek().Type]
		if infixFunction == nil {
			break
		}

		parser.consume()
		expression = infixFunction(context, expression)
	}

	if inferredType, isNever := expression.Type(context).(*types.NeverType); isNever {
		if len(inferredType.Message) > 0 {
			parser.error(currentToken, "Invalid expression (%s)", inferredType.Message)
		}
	}

	return expression
}

/** prefix expressions **/

func (parser *Parser) parsePrefixExpression(context *Context) Expression {
	currentToken := parser.consume()
	return &PrefixExpression{
		PrefixToken: currentToken,
		Operator:    currentToken.Type,
		Expression:  parser.parseExpression(context, Prefix),
	}
}

func (parser *Parser) parseIdentifier(*Context) Expression {
	return &Identifier{IdentToken: parser.current(), Value: parser.current().Literal}
}

func (parser *Parser) parseStringLiteral(*Context) Expression {
	currentToken := parser.current()
	return &StringLiteral{LiteralToken: currentToken, Value: currentToken.Literal}
}

func (parser *Parser) parseBooleanLiteral(*Context) Expression {
	currentToken := parser.current()
	return &BooleanLiteral{LiteralToken: currentToken, Value: currentToken.Type == token.True}
}

func (parser *Parser) parseNullLiteral(*Context) Expression {
	return &NullLiteral{LiteralToken: parser.current()}
}

func (parser *Parser) parseVoidLiteral(*Context) Expression {
	return &VoidLiteral{LiteralToken: parser.current()}
}

func (parser *Parser) parseIntegerLiteral(*Context) Expression {
	currentToken := parser.current()
	literal := &IntegerLiteral{LiteralToken: currentToken}

	value, err := strconv.ParseInt(currentToken.Literal, 10, 64)
	if err != nil {
		parser.error(currentToken, "Integer out of bounds")
		return &InvalidExpression{currentToken}
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parseGroupedExpression(context *Context) Expression {
	parser.consume()
	expression := parser.parseExpression(context, Lowest)
	if !parser.assertNext(token.RParen) {
		return &InvalidExpression{parser.current()}
	}
	return expression
}

func (parser *Parser) parseIncrementPrefixExpression(context *Context) Expression {
	operatorToken := parser.consume()
	identExpression := parser.parseExpression(context, Prefix)
	return parser.parseIncrementExpression(operatorToken, identExpression, true)
}

/** infix expressions **/

func (parser *Parser) parseInfixExpression(context *Context, left Expression) Expression {
	currentToken := parser.consume()
	precedence := precedences[currentToken.Type]

	right := parser.parseExpression(context, precedence)

	return &InfixExpression{
		OperatorToken: currentToken,
		Left:          left,
		Operator:      currentToken.Type,
		Right:         right,
	}
}

func (parser *Parser) parseAssignmentExpression(context *Context, left Expression) Expression {
	assignToken := parser.consume()
	right := parser.parseExpression(context, Assignment)

	ident, isIdent := left.(*Identifier)
	if !isIdent {
		erroneousToken := left.Token()
		parser.error(erroneousToken, "Invalid identifier")
		return &InvalidExpression{InvalidToken: assignToken}
	}

	return &AssignmentExpression{
		IdentToken:  ident.IdentToken,
		AssignToken: assignToken,
		Name:        ident,
		Expression:  right,
	}
}

func (parser *Parser) parseCallExpression(context *Context, function Expression) Expression {
	currentToken := parser.consume()
	argumentList := parser.parseArgumentList(context)

	if argumentList == nil {
		return &InvalidExpression{}
	}

	for _, argument := range argumentList {
		if isInvalid(argument) {
			return argument
		}
	}

	return &CallExpression{
		ParenToken: currentToken,
		Function:   function,
		Arguments:  argumentList,
	}
}

func (parser *Parser) parseIncrementInfixExpression(_ *Context, identExpression Expression) Expression {
	operatorToken := parser.current()
	return parser.parseIncrementExpression(operatorToken, identExpression, false)
}

func (parser *Parser) parseMemberAccessExpression(context *Context, left Expression) Expression {
	dotToken := parser.consume()
	leftType := left.Type(context)
	right := parser.parseExpression(NewSubContext(context, leftType), Postfix)

	ident, isIdent := right.(*Identifier)
	if !isIdent {
		erroneousToken := right.Token()
		parser.error(erroneousToken, "Invalid identifier")
		return &InvalidExpression{InvalidToken: dotToken}
	}

	return &MemberAccessExpression{
		DotToken:   dotToken,
		Expression: left,
		Member:     ident,
		ParentType: leftType,
	}
}

/** misc **/

func (parser *Parser) parseIncrementExpression(operatorToken *token.Token, identExpression Expression, pre bool) Expression {

	ident, isIdent := identExpression.(*Identifier)
	if !isIdent {
		erroneousToken := identExpression.Token()
		parser.error(erroneousToken, "Invalid identifier")
		return &InvalidExpression{erroneousToken}
	}

	return &IncrementExpression{
		IdentToken: ident.IdentToken,
		Name:       ident,
		Operator:   operatorToken.Type,
		Pre:        pre,
	}
}

func (parser *Parser) parseArgumentList(context *Context) []Expression {

	arguments := make([]Expression, 0)

	if parser.current().Type == token.RParen {
		return arguments
	}

	for {
		arguments = append(arguments, parser.parseExpression(context, Lowest))
		if parser.peek().Type == token.Comma {
			parser.consume()
			parser.consume()
		} else {
			break
		}
	}

	if !parser.assertNext(token.RParen) {
		return nil
	}
	return arguments
}

func isInvalid(expression Expression) bool {
	_, invalid := expression.(*InvalidExpression)
	return invalid
}
