package token

type Type int

type Token struct {
	Type    Type
	Literal string
	Line    int
	Col     int
	File    *string
}

func New(tokenType Type, literal string, line int, col int, file *string) *Token {
	return &Token{Type: tokenType, Literal: literal, Line: line, Col: col, File: file}
}

const (
	Illegal Type = iota
	EOF

	Ident
	IntLiteral
	StringLiteral

	EQ
	NEQ
	LT
	GT
	LTE
	GTE

	Plus
	Minus
	Slash
	Star

	LogicalAnd
	LogicalOr

	Assign
	Qmark
	Amp
	Bang
	Increment
	Decrement

	Dot
	Comma
	Semi
	Colon
	DoubleColon
	Define

	LParen
	RParen
	LBrace
	RBrace

	Func
	TypeDef
	Iface
	Return
	Let
	Const
	If
	Else
	For
	While

	True
	False
	Null
	Void
)

var Keywords = map[string]Type{
	"fn":     Func,
	"return": Return,
	"let":    Let,
	"const":  Const,
	"true":   True,
	"false":  False,
	"null":   Null,
	"void":   Void,
	"if":     If,
	"else":   Else,
	"for":    For,
	"while":  While,
	"type":   TypeDef,
	"iface":  Iface,
}

func (token Token) ToString() string {
	return token.Type.ToStringHumanReadable()
}

func (tokenType Type) ToString() string {
	return [...]string{
		"ILLEGAL",
		"EOF",
		"IDENT",
		"INT_LITERAL",
		"STRING_LITERAL",
		"==",
		"!=",
		"<",
		">",
		"<=",
		">=",
		"+",
		"-",
		"/",
		"*",
		"&&",
		"||",
		"=",
		"?",
		"&",
		"!",
		"++",
		"--",
		".",
		",",
		";",
		":",
		"::",
		":=",
		"(",
		")",
		"{",
		"}",
		"FUNC",
		"TYPE",
		"IFACE",
		"RETURN",
		"LET",
		"CONST",
		"IF",
		"ELSE",
		"FOR",
		"WHILE",
		"TRUE",
		"FALSE",
		"NULL",
		"VOID",
	}[tokenType]
}

func (tokenType Type) ToStringHumanReadable() string {
	return [...]string{
		"illegal token",
		"EOF",
		"identifier",
		"integer literal",
		"string literal",
		"'=='",
		"'!='",
		"'<'",
		"'>'",
		"'<='",
		"'>='",
		"'+'",
		"'-'",
		"'/'",
		"'*'",
		"'&&'",
		"'||'",
		"'='",
		"'?'",
		"'&'",
		"'!'",
		"'++'",
		"'--'",
		"'.'",
		"','",
		"';'",
		"':'",
		"'::'",
		"':='",
		"'('",
		"')'",
		"'{'",
		"'}'",
		"'fn'",
		"'type'",
		"'iface'",
		"'return'",
		"'let'",
		"'const'",
		"'if'",
		"'else'",
		"'for'",
		"'while'",
		"'true'",
		"'false'",
		"'null'",
		"'void'",
	}[tokenType]
}
