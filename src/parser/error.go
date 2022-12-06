package parser

import (
	"bananascript/src/token"
	"fmt"
	"github.com/gookit/color"
)

type Error struct {
	Token   *token.Token
	Message string
}

func NewError(token *token.Token, messageFormat string, args ...interface{}) *Error {
	return &Error{Token: token, Message: fmt.Sprintf(messageFormat, args...)}
}

func (error *Error) PrettyPrint(withSource bool) string {
	result := color.FgRed.Sprintf("Error: %s", error.Message)
	if withSource {
		result += "\n\tin "
		if error.Token.File != nil {
			result += *error.Token.File + ":"
		}
		result += fmt.Sprintf("%d:%d", error.Token.Line, error.Token.Col)
	}
	return result
}
