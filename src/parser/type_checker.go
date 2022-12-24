package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
)

func (parser *Parser) getExpressionType(expression Expression, context *types.Context) types.Type {
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
		return &types.String{}
	case *IntegerLiteral:
		return &types.Int{}
	case *BooleanLiteral:
		return &types.Bool{}
	case *NullLiteral:
		return &types.Null{}
	case *VoidLiteral:
		return &types.Void{}
	case *InvalidExpression:
		return &types.Never{}
	}
	parser.error(expression.Token(), "Unknown expression: %T", expression)
	return &types.Never{}
}

func (parser *Parser) getIdentifierType(identifier *Identifier, context *types.Context) types.Type {
	theType, ok := context.GetMemberType(identifier.Value)
	if !ok {
		parser.error(identifier.IdentToken, "Cannot resolve reference to '%s'", identifier.Value)
		return &types.Never{}
	}
	return theType
}

func (parser *Parser) getPrefixExpressionType(prefixExpression *PrefixExpression, context *types.Context) types.Type {

	currentType := parser.getExpressionType(prefixExpression.Expression, context)
	if isNever(currentType) {
		return &types.Never{}
	}

	switch prefixExpression.Operator {
	case token.Bang:
		return &types.Bool{}
	case token.Minus:
		switch currentType.(type) {
		case *types.Int:
			return &types.Int{}
		}
	}

	parser.error(prefixExpression.PrefixToken, "Type mismatch: %s%s", prefixExpression.Operator.ToString(),
		currentType.ToString())
	return &types.Never{}
}

func (parser *Parser) getInfixExpressionType(infixExpression *InfixExpression, context *types.Context) types.Type {
	leftType, rightType := parser.getExpressionType(infixExpression.Left, context), parser.getExpressionType(infixExpression.Right, context)
	if isNever(leftType) || isNever(rightType) {
		return &types.Never{}
	}

	switch infixExpression.Operator {
	case token.EQ, token.NEQ, token.LogicalOr, token.LogicalAnd:
		return &types.Bool{}
	case token.LT, token.GT, token.LTE, token.GTE:
		if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.Bool{}
		}
	case token.Plus:
		if leftType.ToString() == types.TypeString || rightType.ToString() == types.TypeString {
			return &types.String{}
		} else if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.Int{}
		}
	case token.Minus, token.Slash, token.Star:
		if leftType.ToString() == types.TypeInt && rightType.ToString() == types.TypeInt {
			return &types.Int{}
		}
	}

	parser.error(infixExpression.OperatorToken, "Type mismatch: %s %s %s", leftType.ToString(),
		infixExpression.Operator.ToString(), rightType.ToString())
	return &types.Never{}
}

func (parser *Parser) getAssignmentExpressionType(assignmentExpression *AssignmentExpression, context *types.Context) types.Type {
	leftType, rightType := parser.getExpressionType(assignmentExpression.Name, context), parser.getExpressionType(assignmentExpression.Expression, context)
	if isNever(leftType) || isNever(rightType) {
		return &types.Never{}
	}

	if !leftType.IsAssignable(rightType, context) {
		parser.error(assignmentExpression.AssignToken, "Type '%s' is not assignable to '%s'",
			rightType.ToString(), leftType.ToString())
		return &types.Never{}
	}
	return rightType
}

func (parser *Parser) getCallExpressionType(callExpression *CallExpression, context *types.Context) types.Type {
	functionType := parser.getExpressionType(callExpression.Function, context)

	switch functionType := functionType.(type) {
	case *types.Never:
		return &types.Never{}
	case *types.Function:
		if len(functionType.ParameterTypes) == len(callExpression.Arguments) {
			for i, parameterType := range functionType.ParameterTypes {
				if isNever(parameterType) {
					continue
				}
				argumentType := parser.getExpressionType(callExpression.Arguments[i], context)
				if !isNever(argumentType) && !parameterType.IsAssignable(argumentType, context) {
					parser.error(callExpression.Arguments[i].Token(), "Type '%s' is not assignable to '%s'",
						argumentType.ToString(), parameterType.ToString())
				}
			}
		} else {
			parser.error(callExpression.ParenToken, "Mismatching amount of arguments (%d vs %d)",
				len(callExpression.Arguments), len(functionType.ParameterTypes))
		}
		return functionType.ReturnType
	default:
		parser.error(callExpression.ParenToken, "Cannot call '%s'", functionType.ToString())
		return &types.Never{}
	}
}

func (parser *Parser) getIncrementExpressionType(incrementExpression *IncrementExpression, context *types.Context) types.Type {
	identType := parser.getExpressionType(incrementExpression.Name, context)
	switch identType.(type) {
	case *types.Never, *types.Int:
		return identType
	default:
		parser.error(incrementExpression.OperatorToken, "Unknown operator: %s%s",
			incrementExpression.Operator.ToString(), identType.ToString())
		return &types.Never{}
	}
}

func (parser *Parser) getMemberAccessExpressionType(memberAccessExpression *MemberAccessExpression, _ *types.Context) types.Type {
	if isNever(memberAccessExpression.ParentType) {
		return memberAccessExpression.ParentType
	}
	return memberAccessExpression.MemberType
}

func isNever(theType types.Type) bool {
	_, isNever := theType.(*types.Never)
	return isNever
}
