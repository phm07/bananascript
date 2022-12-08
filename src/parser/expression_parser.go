package parser

import (
	"bananascript/src/token"
	"strconv"
)

type ExpressionPrecedence int

const (
	ExpressionLowest ExpressionPrecedence = iota
	ExpressionAssignment
	ExpressionLogicalOr
	ExpressionLogicalAnd
	ExpressionEquals
	ExpressionRelation
	ExpressionSum
	ExpressionProduct
	ExpressionPrefix
	ExpressionPostfix
)

var expressionPrecedences = map[token.Type]ExpressionPrecedence{
	token.Assign:     ExpressionAssignment,
	token.LogicalOr:  ExpressionLogicalOr,
	token.LogicalAnd: ExpressionLogicalAnd,
	token.EQ:         ExpressionEquals,
	token.NEQ:        ExpressionEquals,
	token.LT:         ExpressionRelation,
	token.GT:         ExpressionRelation,
	token.LTE:        ExpressionRelation,
	token.GTE:        ExpressionRelation,
	token.Plus:       ExpressionSum,
	token.Minus:      ExpressionSum,
	token.Slash:      ExpressionProduct,
	token.Star:       ExpressionProduct,
	token.Increment:  ExpressionPostfix,
	token.Decrement:  ExpressionPostfix,
	token.LParen:     ExpressionPostfix,
	token.Dot:        ExpressionPostfix,
}

var prefixExpressionParseFunctions = make(map[token.Type]func(*Context) Expression)
var infixExpressionParseFunctions = make(map[token.Type]func(*Context, Expression) Expression)

func getExpressionPrecedence(token *token.Token) ExpressionPrecedence {
	if precedence, exists := expressionPrecedences[token.Type]; exists {
		return precedence
	}
	return ExpressionLowest
}

func (parser *Parser) registerExpressionParseFunctions() {
	prefixExpressionParseFunctions[token.Ident] = parser.parseIdentifier
	prefixExpressionParseFunctions[token.IntLiteral] = parser.parseIntegerLiteral
	prefixExpressionParseFunctions[token.StringLiteral] = parser.parseStringLiteral
	prefixExpressionParseFunctions[token.Null] = parser.parseNullLiteral
	prefixExpressionParseFunctions[token.Void] = parser.parseVoidLiteral
	prefixExpressionParseFunctions[token.True] = parser.parseBooleanLiteral
	prefixExpressionParseFunctions[token.False] = parser.parseBooleanLiteral
	prefixExpressionParseFunctions[token.Bang] = parser.parsePrefixExpression
	prefixExpressionParseFunctions[token.Minus] = parser.parsePrefixExpression
	prefixExpressionParseFunctions[token.LParen] = parser.parseGroupedExpression
	prefixExpressionParseFunctions[token.Increment] = parser.parseIncrementPrefixExpression
	prefixExpressionParseFunctions[token.Decrement] = parser.parseIncrementPrefixExpression

	infixExpressionParseFunctions[token.Assign] = parser.parseAssignmentExpression
	infixExpressionParseFunctions[token.LogicalOr] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.LogicalAnd] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.EQ] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.NEQ] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.GT] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.LT] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.GTE] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.LTE] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.Plus] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.Minus] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.Slash] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.Star] = parser.parseInfixExpression
	infixExpressionParseFunctions[token.LParen] = parser.parseCallExpression
	infixExpressionParseFunctions[token.Increment] = parser.parseIncrementInfixExpression
	infixExpressionParseFunctions[token.Decrement] = parser.parseIncrementInfixExpression
	infixExpressionParseFunctions[token.Dot] = parser.parseMemberAccessExpression
}

func (parser *Parser) parseExpression(context *Context, precedence ExpressionPrecedence) Expression {
	currentToken := parser.current()
	prefixFunction := prefixExpressionParseFunctions[currentToken.Type]
	if prefixFunction == nil {
		parser.error(currentToken, "Unexpected %s", currentToken.ToString())
		return &InvalidExpression{currentToken}
	}

	expression := prefixFunction(context)

	for parser.peek().Type != token.Semi && precedence < getExpressionPrecedence(parser.peek()) {
		infixFunction := infixExpressionParseFunctions[parser.peek().Type]
		if infixFunction == nil {
			break
		}

		parser.consume()
		expression = infixFunction(context, expression)
	}

	return expression
}

/** prefix expressions **/

func (parser *Parser) parsePrefixExpression(context *Context) Expression {
	currentToken := parser.consume()
	return &PrefixExpression{
		PrefixToken: currentToken,
		Operator:    currentToken.Type,
		Expression:  parser.parseExpression(context, ExpressionPrefix),
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
	expression := parser.parseExpression(context, ExpressionLowest)
	if !parser.assertNext(token.RParen) {
		return &InvalidExpression{parser.current()}
	}
	return expression
}

func (parser *Parser) parseIncrementPrefixExpression(context *Context) Expression {
	operatorToken := parser.consume()
	identExpression := parser.parseExpression(context, ExpressionPrefix)
	return parser.parseIncrementExpression(operatorToken, identExpression, true)
}

/** infix expressions **/

func (parser *Parser) parseInfixExpression(context *Context, left Expression) Expression {
	currentToken := parser.consume()
	precedence := expressionPrecedences[currentToken.Type]

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
	right := parser.parseExpression(context, ExpressionAssignment)

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
	leftType := parser.getExpressionType(left, context)
	right := parser.parseExpression(NewSubContext(context, leftType), ExpressionPostfix)

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
		OperatorToken: operatorToken,
		Name:          ident,
		Operator:      operatorToken.Type,
		Pre:           pre,
	}
}

func (parser *Parser) parseArgumentList(context *Context) []Expression {

	arguments := make([]Expression, 0)

	if parser.current().Type == token.RParen {
		return arguments
	}

	for {
		arguments = append(arguments, parser.parseExpression(context, ExpressionLowest))
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
