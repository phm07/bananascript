package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
	"fmt"
)

func (identifier *Identifier) Type(context *Context) types.Type {
	theType, ok := context.GetType(identifier.Value, nil)
	if !ok {
		if context.parentType != nil {
			return newNever("'%s' is not a member of '%s'", identifier.Value, context.parentType.ToString())
		} else {
			return newNever("Cannot resolve reference to '%s'", identifier.Value)
		}
	}
	return theType
}

func (prefixExpression *PrefixExpression) Type(context *Context) types.Type {

	currentType := prefixExpression.Expression.Type(context)
	if isNever(currentType) {
		return currentType
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

	return newNever("Type mismatch: %s%s", prefixExpression.Operator.ToString(), currentType.ToString())
}

func (infixExpression *InfixExpression) Type(context *Context) types.Type {
	leftType := infixExpression.Left.Type(context)
	if isNever(leftType) {
		return leftType
	}
	rightType := infixExpression.Right.Type(context)
	if isNever(rightType) {
		return rightType
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

	return newNever("Type mismatch: %s %s %s", leftType.ToString(), infixExpression.Operator.ToString(), rightType.ToString())
}

func (assignmentExpression *AssignmentExpression) Type(context *Context) types.Type {

	leftType := assignmentExpression.Name.Type(context)
	if isNever(leftType) {
		return leftType
	}
	rightType := assignmentExpression.Expression.Type(context)
	if isNever(rightType) {
		return rightType
	}

	if !rightType.IsAssignable(leftType) {
		return newNever("Type '%s' is not assignable to '%s'", rightType.ToString(), leftType.ToString())
	}
	return assignmentExpression.Expression.Type(context)
}

func (callExpression *CallExpression) Type(context *Context) types.Type {
	functionType := callExpression.Function.Type(context)

	switch functionType := functionType.(type) {
	case *types.NeverType:
		return functionType
	case *types.FunctionType:
		if len(functionType.ParameterTypes) == len(callExpression.Arguments) {
			for i, parameterType := range functionType.ParameterTypes {
				if isNever(parameterType) {
					return newNever("") // fail silently; error for invalid type already emitted
				}
				argumentType := callExpression.Arguments[i].Type(context)
				if isNever(argumentType) {
					return argumentType
				}
				if !parameterType.IsAssignable(argumentType) {
					return newNever("Type '%s' is not assignable to '%s'", argumentType.ToString(), parameterType.ToString())
				}
			}
			return functionType.ReturnType
		} else {
			return newNever("Mismatching amount of arguments (%d vs %d)", len(callExpression.Arguments),
				len(functionType.ParameterTypes))
		}
	default:
		return newNever("Cannot call %s", functionType.ToString())
	}
}

func (incrementExpression *IncrementExpression) Type(context *Context) types.Type {
	identType := incrementExpression.Name.Type(context)
	switch identType.(type) {
	case *types.NeverType, *types.IntType:
		return identType
	default:
		return newNever("Unknown operator: %s%s", incrementExpression.Operator.ToString(), identType.ToString())
	}
}

func (memberAccessExpression *MemberAccessExpression) Type(context *Context) types.Type {
	if memberType, ok := context.GetType(memberAccessExpression.Member.Value, memberAccessExpression.ParentType); ok {
		return memberType
	} else {
		return newNever("") // member does not exist, cannot resolve reference already caught
	}
}

func (stringLiteral *StringLiteral) Type(*Context) types.Type {
	return &types.StringType{}
}

func (integerLiteral *IntegerLiteral) Type(*Context) types.Type {
	return &types.IntType{}
}

func (booleanLiteral *BooleanLiteral) Type(*Context) types.Type {
	return &types.BoolType{}
}

func (nullLiteral *NullLiteral) Type(*Context) types.Type {
	return &types.NullType{}
}

func (voidLiteral *VoidLiteral) Type(*Context) types.Type {
	return &types.VoidType{}
}

func (invalidExpression *InvalidExpression) Type(*Context) types.Type {
	return &types.NeverType{}
}

func newNever(format string, args ...interface{}) *types.NeverType {
	return &types.NeverType{Message: fmt.Sprintf(format, args...)}
}

func isNever(theType types.Type) bool {
	_, isNever := theType.(*types.NeverType)
	return isNever
}
