package parser

import (
	"bananascript/src/token"
	"fmt"
)

func (identifier *Identifier) Type(context *Context) Type {
	theType, ok := context.GetType(identifier.Value)
	if !ok {
		return newNever("Cannot resolve reference to '%s'", identifier.Value)
	}
	return theType
}

func (prefixExpression *PrefixExpression) Type(context *Context) Type {

	currentType := prefixExpression.Expression.Type(context)
	if isNever(currentType) {
		return currentType
	}

	switch prefixExpression.Operator {
	case token.Bang:
		return &BoolType{}
	case token.Minus:
		switch currentType.(type) {
		case *IntType:
			return &IntType{}
		}
	}

	return newNever("Type mismatch: %s%s", prefixExpression.Operator.ToString(), currentType.ToString())
}

func (infixExpression *InfixExpression) Type(context *Context) Type {
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
		return &BoolType{}
	case token.LT, token.GT, token.LTE, token.GTE:
		if leftType.ToString() == TypeInt && rightType.ToString() == TypeInt {
			return &BoolType{}
		}
	case token.Plus:
		if leftType.ToString() == TypeString || rightType.ToString() == TypeString {
			return &StringType{}
		} else if leftType.ToString() == TypeInt && rightType.ToString() == TypeInt {
			return &IntType{}
		}
	case token.Minus, token.Slash, token.Star:
		if leftType.ToString() == TypeInt && rightType.ToString() == TypeInt {
			return &IntType{}
		}
	}

	return newNever("Type mismatch: %s %s %s", leftType.ToString(), infixExpression.Operator.ToString(), rightType.ToString())
}

func (assignmentExpression *AssignmentExpression) Type(context *Context) Type {

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

func (callExpression *CallExpression) Type(context *Context) Type {
	functionType := callExpression.Function.Type(context)

	switch functionType := functionType.(type) {
	case *NeverType:
		return functionType
	case *FunctionType:
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

func (incrementExpression IncrementExpression) Type(context *Context) Type {
	identType := incrementExpression.Name.Type(context)
	switch identType.(type) {
	case *NeverType, *IntType:
		return identType
	default:
		return newNever("Unknown operator: %s%s", incrementExpression.Operator.ToString(), identType.ToString())
	}
}

func (stringLiteral *StringLiteral) Type(*Context) Type {
	return &StringType{}
}

func (integerLiteral *IntegerLiteral) Type(*Context) Type {
	return &IntType{}
}

func (booleanLiteral *BooleanLiteral) Type(*Context) Type {
	return &BoolType{}
}

func (nullLiteral *NullLiteral) Type(*Context) Type {
	return &NullType{}
}

func (voidLiteral *VoidLiteral) Type(*Context) Type {
	return &VoidType{}
}

func (invalidExpression *InvalidExpression) Type(*Context) Type {
	return &NeverType{}
}

func newNever(format string, args ...interface{}) *NeverType {
	return &NeverType{Message: fmt.Sprintf(format, args...)}
}

func isNever(theType Type) bool {
	_, isNever := theType.(*NeverType)
	return isNever
}
