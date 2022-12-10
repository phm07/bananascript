package errors

import (
	"bananascript/src/token"
	"fmt"
	"github.com/gookit/color"
)

type ParserError struct {
	Line    int
	Col     int
	File    *string
	Message string
}

func New(line int, col int, file *string, messageFormat string, args ...interface{}) *ParserError {
	return &ParserError{Line: line, Col: col, File: file, Message: fmt.Sprintf(messageFormat, args...)}
}

func NewFromToken(token *token.Token, messageFormat string, args ...interface{}) *ParserError {
	return New(token.Line, token.Col, token.File, messageFormat, args...)
}

func (error *ParserError) PrettyPrint(withSource bool) string {
	result := color.FgRed.Sprintf("Error: %s", error.Message)
	if withSource {
		result += "\n\tin "
		if error.File != nil {
			result += *error.File + ":"
		}
		result += fmt.Sprintf("%d:%d", error.Line, error.Col)
	}
	return result
}
