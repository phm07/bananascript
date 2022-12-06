package evaluator

import (
	"bananascript/src/parser"
	"bananascript/src/token"
	"fmt"
	"reflect"
)

func Eval(node parser.Node, environment *Environment) Object {
	switch node := node.(type) {
	case *parser.Program:
		return evalProgram(node, environment)
	case *parser.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *parser.StringLiteral:
		return &StringObject{Value: node.Value}
	case *parser.IntegerLiteral:
		return &IntegerObject{Value: node.Value}
	case *parser.BooleanLiteral:
		return &BooleanObject{Value: node.Value}
	case *parser.Identifier:
		return evalIdentifierExpression(node, environment)
	case *parser.InfixExpression:
		return evalInfixExpression(node, environment)
	case *parser.PrefixExpression:
		return evalPrefixExpression(node, environment)
	case *parser.CallExpression:
		return evalCallExpression(node, environment)
	case *parser.AssignmentExpression:
		return evalAssignmentExpression(node, environment)
	case *parser.LetStatement:
		return evalLetStatement(node, environment)
	case *parser.FunctionDefinitionStatement:
		return evalFunctionDefinitionStatement(node, environment)
	case *parser.ReturnStatement:
		return evalReturnStatement(node, environment)
	case *parser.BlockStatement:
		return evalBlockStatement(node, environment)
	case *parser.IfStatement:
		return evalIfStatement(node, environment)
	case *parser.WhileStatement:
		return evalWhileStatement(node, environment)
	case *parser.IncrementExpression:
		return evalIncrementExpression(node, environment)
	}
	return nil
}

func evalProgram(program *parser.Program, environment *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement, environment)
		switch result := result.(type) {
		case *ErrorObject:
			return result
		}
	}
	return result
}

func evalPrefixExpression(prefixExpression *parser.PrefixExpression, environment *Environment) Object {

	object := Eval(prefixExpression.Expression, environment)
	if isError(object) {
		return object
	}

	switch prefixExpression.Operator {
	case token.Bang:
		return &BooleanObject{Value: !implicitBoolConversion(object)}
	case token.Minus:
		intValue, ok := intConversion(object)
		if ok {
			return &IntegerObject{Value: -intValue}
		}
	}

	return NewError("Unknown prefix operator")
}

func evalInfixExpression(infixExpression *parser.InfixExpression, environment *Environment) Object {

	leftObject := Eval(infixExpression.Left, environment)
	if isError(leftObject) {
		return leftObject
	}
	rightObject := Eval(infixExpression.Right, environment)
	if isError(rightObject) {
		return rightObject
	}

	switch infixExpression.Operator {
	case token.EQ:
		return &BooleanObject{Value: evalEquals(leftObject, rightObject)}
	case token.NEQ:
		return &BooleanObject{Value: !evalEquals(leftObject, rightObject)}
	case token.LT:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &BooleanObject{Value: left < right} })
	case token.GT:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &BooleanObject{Value: left > right} })
	case token.LTE:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &BooleanObject{Value: left <= right} })
	case token.GTE:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &BooleanObject{Value: left >= right} })
	case token.Plus:
		_, leftIsString := leftObject.(*StringObject)
		_, rightIsString := rightObject.(*StringObject)
		if leftIsString || rightIsString {
			return &StringObject{Value: leftObject.ToString() + rightObject.ToString()}
		}
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &IntegerObject{Value: left + right} })
	case token.Minus:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &IntegerObject{Value: left - right} })
	case token.Slash:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &IntegerObject{Value: left / right} })
	case token.Star:
		return evalIntegerInfix(leftObject, rightObject,
			func(left int64, right int64) Object { return &IntegerObject{Value: left * right} })
	default:
		return NewError("Unknown infix operator")
	}
}

func evalEquals(left Object, right Object) bool {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false
	}
	switch left := left.(type) {
	case *BooleanObject:
		return left.Value == right.(*BooleanObject).Value
	case *IntegerObject:
		return left.Value == right.(*IntegerObject).Value
	case *StringObject:
		return left.Value == right.(*StringObject).Value
	default:
		return left == right
	}
}

func evalIntegerInfix(left Object, right Object, constructor func(left int64, right int64) Object) Object {
	leftInt, leftOk := intConversion(left)
	rightInt, rightOk := intConversion(right)
	if !rightOk || !leftOk {
		return NewError("Implicit conversion to int not possible")
	}
	return constructor(leftInt, rightInt)
}

func evalAssignmentExpression(assignmentExpression *parser.AssignmentExpression, environment *Environment) Object {

	object := Eval(assignmentExpression.Expression, environment)
	if isError(object) {
		return object
	}

	name := assignmentExpression.Name.Value
	if object, ok := environment.Assign(name, object); ok {
		return object
	} else {
		return NewError("Cannot resolve variable")
	}
}

func evalCallExpression(callExpression *parser.CallExpression, environment *Environment) Object {
	function := Eval(callExpression.Function, environment)
	switch function := function.(type) {
	case *ErrorObject:
		return function
	case *FunctionObject:
		if len(callExpression.Arguments) != len(function.Parameters) {
			return NewError("Mismatching number of arguments")
		}
		argumentObjects := make([]Object, 0)
		for _, argument := range callExpression.Arguments {
			argumentObjects = append(argumentObjects, Eval(argument, environment))
		}
		returned := function.Execute(argumentObjects)
		switch returned := returned.(type) {
		case *ReturnObject:
			return returned.Object
		default:
			return returned
		}
	default:
		return NewError("Cannot call non-function")
	}
}

func evalIdentifierExpression(identifier *parser.Identifier, environment *Environment) Object {
	if object, exists := environment.Get(identifier.Value); exists {
		return object
	} else {
		return NewError("Cannot resolve identifier")
	}
}

func evalLetStatement(letStatement *parser.LetStatement, environment *Environment) Object {

	object := Eval(letStatement.Value, environment)
	if isError(object) {
		return object
	}

	name := letStatement.Name.Value
	environment.Define(name, object)
	return nil
}

func evalFunctionDefinitionStatement(funcStatement *parser.FunctionDefinitionStatement, environment *Environment) Object {

	name := funcStatement.Name.Value
	if _, exists := environment.GetInThisScope(name); exists {
		return NewError("Cannot re-declare function")
	}

	identifiers := make([]*parser.Identifier, 0)
	for _, parameter := range funcStatement.Parameters {
		identifiers = append(identifiers, parameter.Name)
	}

	object := &FunctionObject{
		Parameters: identifiers,
		Execute: func(arguments []Object) Object {
			newEnvironment := ExtendEnvironment(environment)
			for i, argument := range arguments {
				name := identifiers[i].Value
				newEnvironment.Define(name, argument)
			}
			return Eval(funcStatement.Body, newEnvironment)
		},
	}
	environment.Define(name, object)
	return nil
}

func evalReturnStatement(returnStatement *parser.ReturnStatement, environment *Environment) Object {
	object := Eval(returnStatement.Expression, environment)
	if isError(object) {
		return object
	}
	return &ReturnObject{Object: object}
}

func evalBlockStatement(blockStatement *parser.BlockStatement, environment *Environment) Object {
	newEnvironment := ExtendEnvironment(environment)

	for _, statement := range blockStatement.Statements {
		object := Eval(statement, newEnvironment)
		if object != nil {
			switch object := object.(type) {
			case *ErrorObject, *ReturnObject:
				return object
			default:
				continue
			}
		}
	}

	return nil
}

func evalIfStatement(ifStatement *parser.IfStatement, environment *Environment) Object {
	condition := Eval(ifStatement.Condition, environment)
	if isError(condition) {
		return condition
	}
	if implicitBoolConversion(condition) {
		return Eval(ifStatement.Statement, ExtendEnvironment(environment))
	} else {
		return Eval(ifStatement.Alternative, ExtendEnvironment(environment))
	}
}

func evalWhileStatement(whileStatement *parser.WhileStatement, environment *Environment) Object {
	for {
		condition := Eval(whileStatement.Condition, environment)
		if isError(condition) {
			return condition
		}
		if !implicitBoolConversion(condition) {
			return nil
		}
		object := Eval(whileStatement.Statement, ExtendEnvironment(environment))
		switch object := object.(type) {
		case *ErrorObject, *ReturnObject:
			return object
		default:
			continue
		}
	}
}

func evalIncrementExpression(incrementExpression *parser.IncrementExpression, environment *Environment) Object {

	object, exists := environment.Get(incrementExpression.Name.Value)
	if !exists {
		return NewError("Cannot resolve identifier")
	}

	if object, ok := object.(*IntegerObject); ok {
		oldValue := object.Value
		if incrementExpression.Operator == token.Increment {
			object.Value++
		} else {
			object.Value--
		}
		if incrementExpression.Pre {
			return object
		} else {
			return &IntegerObject{Value: oldValue}
		}
	} else {
		return NewError("Cannot increment non-int")
	}
}

func implicitBoolConversion(object Object) bool {
	switch object := object.(type) {
	case *BooleanObject:
		return object.Value
	case *IntegerObject:
		return object.Value != 0
	case *StringObject:
		return len(object.Value) != 0
	default:
		return true
	}
}

func intConversion(object Object) (int64, bool) {
	switch object := object.(type) {
	case *IntegerObject:
		return object.Value, true
	default:
		return 0, false
	}
}

func NewError(format string, args ...interface{}) *ErrorObject {
	return &ErrorObject{Message: fmt.Sprintf(format, args...)}
}

func isError(object Object) bool {
	_, isError := object.(*ErrorObject)
	return isError
}
