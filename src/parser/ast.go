package parser

import (
	"bananascript/src/token"
	"bananascript/src/types"
	"fmt"
	"strconv"
)

type Node interface {
	ToString() string
	Token() *token.Token
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Statements []Statement
}

func (program *Program) Token() *token.Token {
	return program.Statements[0].Token()
}

func (program *Program) ToString() string {
	result := ""
	for _, statement := range program.Statements {
		result += statement.ToString() + "\n"
	}
	return result
}

type InvalidExpression struct {
	InvalidToken *token.Token
}

func (invalidExpression *InvalidExpression) Token() *token.Token {
	return invalidExpression.InvalidToken
}

func (invalidExpression *InvalidExpression) ToString() string {
	return ""
}

type Identifier struct {
	IdentToken *token.Token
	Value      string
}

func (identifier *Identifier) Token() *token.Token {
	return identifier.IdentToken
}

func (identifier *Identifier) ToString() string {
	return identifier.Value
}

type Parameter struct {
	Token *token.Token
	Name  *Identifier
	Type  types.Type
}

func (parameter *Parameter) ToString() string {
	return parameter.Name.Value + ": " + parameter.Type.ToString()
}

type ExpressionStatement struct {
	FirstToken *token.Token
	Expression Expression
}

func (expressionStatement *ExpressionStatement) Token() *token.Token {
	return expressionStatement.FirstToken
}

func (expressionStatement *ExpressionStatement) ToString() string {
	return expressionStatement.Expression.ToString()
}

type PrefixExpression struct {
	PrefixToken *token.Token
	Operator    token.Type
	Expression  Expression
}

func (prefixExpression *PrefixExpression) Token() *token.Token {
	return prefixExpression.PrefixToken
}

func (prefixExpression *PrefixExpression) ToString() string {
	return "(" + prefixExpression.Operator.ToString() + prefixExpression.Expression.ToString() + ")"
}

type InfixExpression struct {
	OperatorToken *token.Token
	Left          Expression
	Operator      token.Type
	Right         Expression
}

func (infixExpression *InfixExpression) Token() *token.Token {
	return infixExpression.Left.Token()
}

func (infixExpression *InfixExpression) ToString() string {
	return "(" + infixExpression.Left.ToString() + " " + infixExpression.Operator.ToString() + " " +
		infixExpression.Right.ToString() + ")"
}

type AssignmentExpression struct {
	IdentToken  *token.Token
	AssignToken *token.Token
	Name        *Identifier
	Expression  Expression
}

func (assignmentExpression *AssignmentExpression) Token() *token.Token {
	return assignmentExpression.IdentToken
}

func (assignmentExpression *AssignmentExpression) ToString() string {
	return "(" + assignmentExpression.Name.Value + " = " + assignmentExpression.Expression.ToString() + ")"
}

type CallExpression struct {
	ParenToken *token.Token
	Function   Expression
	Arguments  []Expression
}

func (callExpression *CallExpression) Token() *token.Token {
	return callExpression.ParenToken
}

func (callExpression *CallExpression) ToString() string {
	expressionString := "(" + callExpression.Function.ToString() + "("
	for i, argument := range callExpression.Arguments {
		if i > 0 {
			expressionString += ", "
		}
		expressionString += argument.ToString()
	}
	return expressionString + "))"
}

type StringLiteral struct {
	LiteralToken *token.Token
	Value        string
}

func (stringLiteral *StringLiteral) Token() *token.Token {
	return stringLiteral.LiteralToken
}

func (stringLiteral *StringLiteral) ToString() string {
	return stringLiteral.Value
}

type IntegerLiteral struct {
	LiteralToken *token.Token
	Value        int64
}

func (integerLiteral *IntegerLiteral) Token() *token.Token {
	return integerLiteral.LiteralToken
}

func (integerLiteral *IntegerLiteral) ToString() string {
	return strconv.FormatInt(integerLiteral.Value, 10)
}

type BooleanLiteral struct {
	LiteralToken *token.Token
	Value        bool
}

func (booleanLiteral *BooleanLiteral) Token() *token.Token {
	return booleanLiteral.LiteralToken
}

func (booleanLiteral *BooleanLiteral) ToString() string {
	return strconv.FormatBool(booleanLiteral.Value)
}

type NullLiteral struct {
	LiteralToken *token.Token
}

func (nullLiteral *NullLiteral) Token() *token.Token {
	return nullLiteral.LiteralToken
}

func (nullLiteral *NullLiteral) ToString() string {
	return "null"
}

type VoidLiteral struct {
	LiteralToken *token.Token
}

func (voidLiteral *VoidLiteral) Token() *token.Token {
	return voidLiteral.LiteralToken
}

func (voidLiteral *VoidLiteral) ToString() string {
	return "void"
}

type FunctionDefinitionStatement struct {
	FuncToken  *token.Token
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
	ThisType   types.Type
	ReturnType types.Type
}

func (funcStatement *FunctionDefinitionStatement) Token() *token.Token {
	return funcStatement.FuncToken
}

func (funcStatement *FunctionDefinitionStatement) ToString() string {
	result := "fn " + funcStatement.Name.Value + "("
	for i, parameter := range funcStatement.Parameters {
		if i > 0 {
			result += ", "
		}
		result += parameter.ToString()
	}
	return result + funcStatement.Body.ToString()
}

type LetStatement struct {
	LetToken *token.Token
	Name     *Identifier
	Type     types.Type
	Value    Expression
}

func (letStatement *LetStatement) Token() *token.Token {
	return letStatement.LetToken
}

func (letStatement *LetStatement) ToString() string {
	return fmt.Sprintf("let %s: %s = %s;", letStatement.Name.Value, letStatement.Type.ToString(), letStatement.Value.ToString())
}

type ReturnStatement struct {
	ReturnToken *token.Token
	Expression  Expression
}

func (returnStatement *ReturnStatement) Token() *token.Token {
	return returnStatement.ReturnToken
}

func (returnStatement *ReturnStatement) ToString() string {
	return "return " + returnStatement.Expression.ToString() + ";"
}

type BlockStatement struct {
	LBraceToken *token.Token
	RBraceToken *token.Token
	Statements  []Statement
}

func (blockStatement *BlockStatement) Token() *token.Token {
	return blockStatement.LBraceToken
}

func (blockStatement *BlockStatement) ToString() string {
	result := "{"
	for _, statement := range blockStatement.Statements {
		result += "\n    " + statement.ToString()
	}
	return result + "\n}"
}

type IfStatement struct {
	IfToken     *token.Token
	Condition   Expression
	Statement   Statement
	Alternative Statement
}

func (ifStatement *IfStatement) Token() *token.Token {
	return ifStatement.IfToken
}

func (ifStatement *IfStatement) ToString() string {
	result := "if " + ifStatement.Condition.ToString() + " " + ifStatement.Statement.ToString()
	if ifStatement.Alternative != nil {
		result += " else " + ifStatement.Alternative.ToString()
	}
	return result
}

type WhileStatement struct {
	WhileToken *token.Token
	Condition  Expression
	Statement  Statement
}

func (whileStatement *WhileStatement) Token() *token.Token {
	return whileStatement.WhileToken
}

func (whileStatement *WhileStatement) ToString() string {
	return "while " + whileStatement.Condition.ToString() + " " + whileStatement.Statement.ToString()
}

type IncrementExpression struct {
	OperatorToken *token.Token
	Operator      token.Type
	Name          *Identifier
	Pre           bool
}

func (incrementExpression *IncrementExpression) Token() *token.Token {
	return incrementExpression.OperatorToken
}

func (incrementExpression *IncrementExpression) ToString() string {
	result := ""
	if incrementExpression.Pre {
		result += incrementExpression.Operator.ToString()
	}
	result += incrementExpression.Name.Value
	if !incrementExpression.Pre {
		result += incrementExpression.Operator.ToString()
	}
	return result
}

type MemberAccessExpression struct {
	DotToken   *token.Token
	Expression Expression
	Member     *Identifier
	ParentType types.Type
}

func (memberAccessExpression *MemberAccessExpression) Token() *token.Token {
	return memberAccessExpression.DotToken
}

func (memberAccessExpression *MemberAccessExpression) ToString() string {
	return memberAccessExpression.Expression.ToString() + "." + memberAccessExpression.Member.Value
}

type TypeDefinitionStatement struct {
	IdentToken *token.Token
	Name       *Identifier
	Type       types.Type
}

func (typeDefinitionStatement *TypeDefinitionStatement) Token() *token.Token {
	return typeDefinitionStatement.IdentToken
}

func (typeDefinitionStatement *TypeDefinitionStatement) ToString() string {
	return fmt.Sprintf("type %s := %s;", typeDefinitionStatement.Name.Value, typeDefinitionStatement.Type.ToString())
}
