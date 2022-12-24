package lexer

import (
	"bananascript/src/errors"
	"bananascript/src/token"
	"os"
	"path/filepath"
	"strconv"
)

type Lexer struct {
	Errors         []*errors.ParserError
	input          []rune
	position       int
	line           int
	col            int
	filePath       *string
	lastWasIllegal bool
}

func FromFile(fileName string) (*Lexer, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	filePathAbsolute, err := filepath.Abs(fileName)
	if err != nil {
		return nil, err
	}

	input := string(bytes)
	return &Lexer{input: []rune(input), line: 1, filePath: &filePathAbsolute, Errors: make([]*errors.ParserError, 0)}, nil
}

func FromCode(input string) *Lexer {
	return &Lexer{input: []rune(input), line: 1, Errors: make([]*errors.ParserError, 0)}
}

func (lexer *Lexer) current() rune {
	if lexer.position < len(lexer.input) {
		return lexer.input[lexer.position]
	} else {
		return 0
	}
}

func (lexer *Lexer) consume() rune {
	ch := lexer.current()
	lexer.position++
	if ch == '\n' {
		lexer.line++
		lexer.col = 0
	} else {
		lexer.col++
	}
	return ch
}

func (lexer *Lexer) peek() rune {
	if lexer.position+1 < len(lexer.input) {
		return lexer.input[lexer.position+1]
	} else {
		return 0
	}
}

func (lexer *Lexer) NextToken() *token.Token {

	lexer.eatWhitespace()
	char := lexer.consume()
	startCol := lexer.col

	switch char {
	case 0:
		return lexer.newToken(token.EOF, "", startCol)
	case '=':
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.EQ, "", startCol)
		}
		return lexer.newToken(token.Assign, "", startCol)
	case '+':
		if lexer.current() == '+' {
			lexer.consume()
			return lexer.newToken(token.Increment, "", startCol)
		}
		return lexer.newToken(token.Plus, "", startCol)
	case '-':
		if lexer.current() == '-' {
			lexer.consume()
			return lexer.newToken(token.Decrement, "", startCol)
		}
		return lexer.newToken(token.Minus, "", startCol)
	case '/':
		if lexer.current() == '/' {
			lexer.eatLine()
			return lexer.NextToken()
		} else if lexer.current() == '*' {
			lexer.eatComment()
			return lexer.NextToken()
		}
		return lexer.newToken(token.Slash, "", startCol)
	case '*':
		return lexer.newToken(token.Star, "", startCol)
	case '<':
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.LTE, "", startCol)
		}
		return lexer.newToken(token.LT, "", startCol)
	case '>':
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.GTE, "", startCol)
		}
		return lexer.newToken(token.GT, "", startCol)
	case '?':
		return lexer.newToken(token.Qmark, "", startCol)
	case '&':
		if lexer.current() == '&' {
			lexer.consume()
			return lexer.newToken(token.LogicalAnd, "", startCol)
		}
		return lexer.newToken(token.Amp, "", startCol)
	case '|':
		if lexer.current() == '|' {
			lexer.consume()
			return lexer.newToken(token.LogicalOr, "", startCol)
		}
	case '!':
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.NEQ, "", startCol)
		}
		return lexer.newToken(token.Bang, "", startCol)
	case '.':
		return lexer.newToken(token.Dot, "", startCol)
	case ',':
		return lexer.newToken(token.Comma, "", startCol)
	case ';':
		return lexer.newToken(token.Semi, "", startCol)
	case ':':
		if lexer.current() == ':' {
			lexer.consume()
			return lexer.newToken(token.DoubleColon, "", startCol)
		}
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.Define, "", startCol)
		}
		return lexer.newToken(token.Colon, "", startCol)
	case '(':
		return lexer.newToken(token.LParen, "", startCol)
	case ')':
		return lexer.newToken(token.RParen, "", startCol)
	case '{':
		return lexer.newToken(token.LBrace, "", startCol)
	case '}':
		return lexer.newToken(token.RBrace, "", startCol)
	case '"':
		return lexer.parseString(startCol)
	}

	if isIdent(char) {
		start := lexer.position - 1
		for isIdent(lexer.current()) || isDigit(lexer.current()) {
			lexer.consume()
		}
		ident := string(lexer.input[start:lexer.position])
		if tokenType, exists := token.Keywords[ident]; exists {
			return lexer.newToken(tokenType, "", startCol)
		}
		return lexer.newToken(token.Ident, ident, startCol)

	} else if isDigit(char) {
		start := lexer.position - 1
		isFloat := false
		for isDigit(lexer.current()) {
			lexer.consume()
		}
		if lexer.current() == '.' {
			isFloat = true
			lexer.consume()
			for isDigit(lexer.current()) {
				lexer.consume()
			}
		}
		literal := string(lexer.input[start:lexer.position])
		if isFloat {
			return lexer.newToken(token.FloatLiteral, literal, startCol)
		} else {
			return lexer.newToken(token.IntLiteral, literal, startCol)
		}
	} else {
		if !lexer.lastWasIllegal {
			lexer.error(startCol, "Illegal token")
		}
		return lexer.newToken(token.Illegal, string(char), startCol)
	}
}

func (lexer *Lexer) parseString(stringStartCol int) *token.Token {
	literal := ""

parseChar:
	for {
		current := lexer.consume()
		startCol := lexer.col
		if current == '"' || current == '\n' || current == 0 {
			if current == '\n' || current == 0 {
				lexer.error(startCol+1, "Unclosed string literal")
			}
			return lexer.newToken(token.StringLiteral, literal, stringStartCol)
		}
		toAdd := string(current)
		if current == '\\' {
			switch next := lexer.consume(); next {
			case '\\':
				toAdd = "\\"
			case '"':
				toAdd = "\""
			case '\'':
				toAdd = "'"
			case 'a':
				toAdd = "\a"
			case 'b':
				toAdd = "\b"
			case 'f':
				toAdd = "\f"
			case 'n':
				toAdd = "\n"
			case 'r':
				toAdd = "\r"
			case 't':
				toAdd = "\t"
			case 'v':
				toAdd = "\v"
			case 'x':
				hex := ""
				for {
					current := lexer.current()
					if isHex(current) {
						hex += string(current)
						lexer.consume()
					} else {
						break
					}
				}
				value, err := strconv.ParseInt(hex, 16, 64)
				if err != nil {
					lexer.error(startCol, "Invalid hexadecimal (%s)", hex)
				} else {
					toAdd = string(rune(value))
				}
			case 'u', 'U':
				nDigits := 4
				if next == 'U' {
					nDigits = 8
				}
				hex := ""
				for i := 0; i < nDigits; i++ {
					current := lexer.current()
					if !isHex(current) {
						lexer.error(startCol, "Invalid unicode sequence (%s)", hex)
						continue parseChar
					}
					lexer.consume()
					hex += string(current)
				}
				value, err := strconv.ParseInt(hex, 16, nDigits*8)
				if err != nil {
					lexer.error(startCol, "Invalid unicode sequence (%s)", hex)
				} else {
					toAdd = string(rune(value))
				}
			default:
				if next >= '0' && next <= '8' {
					octal := string(next)
					for i := 0; i < 2; i++ {
						current := lexer.current()
						if current >= '0' && current <= '8' {
							octal += string(current)
							lexer.consume()
						} else {
							break
						}
					}
					value, err := strconv.ParseInt(octal, 8, 16)
					if err != nil {
						lexer.error(startCol, "Invalid octal (%s)", octal)
					} else {
						toAdd = string(rune(value))
					}
				}
				lexer.error(startCol, "Invalid escape sequence")
			}
		}
		literal += toAdd
	}
}

func (lexer *Lexer) eatWhitespace() {
	for isWhitespace(lexer.current()) {
		lexer.consume()
	}
}

func (lexer *Lexer) eatLine() {
	for lexer.current() != 0 && lexer.current() != '\n' {
		lexer.consume()
	}
}

func (lexer *Lexer) eatComment() {
	for lexer.current() != 0 && (lexer.current() != '*' || lexer.peek() != '/') {
		lexer.consume()
	}
	lexer.consume() // *
	lexer.consume() // /
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\r' || char == '\v' || char == '\f' || char == '\n'
}

func isIdent(char rune) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_'
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isHex(char rune) bool {
	return (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')
}

func (lexer *Lexer) error(startCol int, messageFormat string, args ...interface{}) {
	lexer.Errors = append(lexer.Errors, errors.New(lexer.line, startCol, lexer.filePath, messageFormat, args...))
}

func (lexer *Lexer) newToken(tokenType token.Type, literal string, startCol int) *token.Token {
	lexer.lastWasIllegal = tokenType == token.Illegal
	return token.New(tokenType, literal, lexer.line, startCol, lexer.filePath)
}
