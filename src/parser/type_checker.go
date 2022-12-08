package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
)

func (parser *Parser) getExpressionType(expression Expression, context *Context) types.Type {
	switch expression := expression.(type) {
	case *Identifier:
		return parser.getIdentifierType(expression, context)
	case *PrefixExpression:
		return parser.getPrefixExpressionType(expression, context)
	case *InfixExpression:
		return parser.getInfixExpressionType(expression, context)
	case *AssignmentExpression:
		return parser.getAssignmentExpressionType(expression, context)
	case *CallExpression:
		return parser.getCallExpressionType(expression, context)
	case *IncrementExpression:
		return parser.getIncrementExpressionType(expression, context)
	case *MemberAccessExpression:
		return parser.getMemberAccessExpressionType(expression, context)
	case *StringLiteral:
		return &types.StringType{}
	case *IntegerLiteral:
		return &types.IntType{}
	case *BooleanLiteral:
		return &types.BoolType{}
	case *NullLiteral:
		return &types.NullType{}
	case *VoidLiteral:
		return &types.VoidType{}
	case *InvalidExpression:
		return &types.NeverType{}
	}
	parser.error(expression.Token(), "Unknown expression: %T", expression)
	return &types.NeverType{}
}

func (parser *Parser) getIdentifierType(identifier *Identifier, context *Context) types.Type {
	theType, ok := context.GetType(identifier.Value, nil)
	if !ok {
		parser.error(identifier.IdentToken, "Cannot resolve reference to '%s'", identifier.Value)
		return &types.NeverType{}
	}
	return theType
}

func (parser *Parser) getPrefixExpressionType(prefixExpression *PrefixExpression, context *Context) types.Type {

	currentType := parser.getExpressionType(prefixExpression.Expression, context)
	if isNever(currentType) {
		return &types.NeverType{}
	}

	switch prefixExpression.Operator {
	case token.Bang:
		return &types.BoolType{}
	case token.Minus:
		switch currentType.(type) {
		case *types.IntType:
			return &types.IntType{}
		}
	}

	parser.error(prefixExpression.PrefixToken, "Type mismatch: %s%s", prefixExpression.Operator.ToString(),
		currentType.ToString())
	return &types.NeverType{}
}

func (parser *Parser) getInfixExpressionType(infixExpression *InfixExpression, context *Context) types.Type {
	leftType, rightType := parser.getExpressionType(infixExpression.Left, context), parser.getExpressionType(infixExpression.Right, context)
	if isNever(leftType) || isNever(rightType) {
		return &types.NeverType{}
	}

	switch infixExpression.Operator {
	case token.EQ, token.NEQ, token.LogicalOr, token.LogicalAnd:
		return &types.BoolType{}
	case token.LT, token.GT, token.LTE, token.GTE:
		if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.BoolType{}
		}
	case token.Plus:
		if leftType.ToString() == types.TypeString || rightType.ToString() == types.TypeString {
			return &types.StringType{}
		} else if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.IntType{}
		}
	case token.Minus, token.Slash, token.Star:
		if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.IntType{}
		}
	}

	parser.error(infixExpression.OperatorToken, "Type mismatch: %s %s %s", leftType.ToString(),
		infixExpression.Operator.ToString(), rightType.ToString())
	return &types.NeverType{}
}

func (parser *Parser) getAssignmentExpressionType(assignmentExpression *AssignmentExpression, context *Context) types.Type {
	leftType, rightType := parser.getExpressionType(assignmentExpression.Name, context), parser.getExpressionType(assignmentExpression.Expression, context)
	if isNever(leftType) || isNever(rightType) {
		return &types.NeverType{}
	}

	if !rightType.IsAssignable(leftType) {
		parser.error(assignmentExpression.AssignToken, "Type '%s' is not assignable to '%s'",
			rightType.ToString(), leftType.ToString())
		return &types.NeverType{}
	}
	return rightType
}

func (parser *Parser) getCallExpressionType(callExpression *CallExpression, context *Context) types.Type {
	functionType := parser.getExpressionType(callExpression.Function, context)

	switch functionType := functionType.(type) {
	case *types.NeverType:
		return &types.NeverType{}
	case *types.FunctionType:
		if len(functionType.ParameterTypes) == len(callExpression.Arguments) {
			for i, parameterType := range functionType.ParameterTypes {
				if isNever(parameterType) {
					continue
				}
				argumentType := parser.getExpressionType(callExpression.Arguments[i], context)
				if !isNever(argumentType) && !parameterType.IsAssignable(argumentType) {
					parser.error(callExpression.Arguments[i].Token(), "Type '%s' is not assignable to '%s'",
						argumentType.ToString(), parameterType.ToString())
				}
			}
		} else {
			parser.error(callExpression.ParenToken, "Mismatching amount of arguments (%d vs %d)", len(callExpression.Arguments))
		}
		return functionType.ReturnType
	default:
		parser.error(callExpression.ParenToken, "Cannot call '%s'", functionType.ToString())
		return &types.NeverType{}
	}
}

func (parser *Parser) getIncrementExpressionType(incrementExpression *IncrementExpression, context *Context) types.Type {
	identType := parser.getExpressionType(incrementExpression.Name, context)
	switch identType.(type) {
	case *types.NeverType, *types.IntType:
		return identType
	default:
		parser.error(incrementExpression.OperatorToken, "Unknown operator: %s%s",
			incrementExpression.Operator.ToString(), identType.ToString())
		return &types.NeverType{}
	}
}

func (parser *Parser) getMemberAccessExpressionType(memberAccessExpression *MemberAccessExpression, context *Context) types.Type {
	if isNever(memberAccessExpression.ParentType) {
		return memberAccessExpression.ParentType
	}
	if memberType, ok := context.GetType(memberAccessExpression.Member.Value, memberAccessExpression.ParentType); ok {
		return memberType
	} else {
		parser.error(memberAccessExpression.DotToken, "Member '%s' does not exist on '%s'",
			memberAccessExpression.Member.Value, memberAccessExpression.ParentType.ToString())
		return &types.NeverType{}
	}
}

func isNever(theType types.Type) bool {
	_, isNever := theType.(*types.NeverType)
	return isNever
}
