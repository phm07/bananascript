package lexer

import (
	"bananascript/src/token"
	"os"
	"path/filepath"
)

type Lexer struct {
	input    string
	position int
	line     int
	col      int
	filePath *string
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
	return &Lexer{input: input, line: 1, filePath: &filePathAbsolute}, nil
}

func FromCode(input string) *Lexer {
	return &Lexer{input: input, line: 1}
}

func (lexer *Lexer) current() byte {
	if lexer.position < len(lexer.input) {
		return lexer.input[lexer.position]
	} else {
		return 0
	}
}

func (lexer *Lexer) consume() byte {
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

func (lexer *Lexer) peek() byte {
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
		return lexer.newToken(token.Amp, "", startCol)
	case '!':
		if lexer.current() == '=' {
			lexer.consume()
			return lexer.newToken(token.NEQ, "", startCol)
		}
		return lexer.newToken(token.Bang, "", startCol)
	case ',':
		return lexer.newToken(token.Comma, "", startCol)
	case ';':
		return lexer.newToken(token.Semi, "", startCol)
	case ':':
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
		start := lexer.position
		for !endsString(lexer.current()) {
			lexer.consume()
		}
		end := lexer.position
		if lexer.consume() == '"' {
			literal := lexer.input[start:end]
			return lexer.newToken(token.StringLiteral, literal, startCol)
		} else {
			literal := lexer.input[(start - 1):end]
			return lexer.newToken(token.Illegal, literal, startCol)
		}
	default:
		if isIdent(char) {
			start := lexer.position - 1
			for isIdent(lexer.current()) || isDigit(lexer.current()) {
				lexer.consume()
			}
			ident := lexer.input[start:lexer.position]
			if tokenType, exists := token.Keywords[ident]; exists {
				return lexer.newToken(tokenType, "", startCol)
			}
			return lexer.newToken(token.Ident, ident, startCol)

		} else if isDigit(char) {
			start := lexer.position - 1
			for isDigit(lexer.current()) {
				lexer.consume()
			}
			integer := lexer.input[start:lexer.position]
			return lexer.newToken(token.IntLiteral, integer, startCol)

		} else {
			return lexer.newToken(token.Illegal, string(char), startCol)
		}
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

func endsString(char byte) bool {
	return char == '"' || char == '\n' || char == 0
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\r' || char == '\v' || char == '\f' || char == '\n'
}

func isIdent(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (lexer *Lexer) newToken(tokenType token.Type, literal string, startCol int) *token.Token {
	return token.New(tokenType, literal, lexer.line, startCol, lexer.filePath)
}
