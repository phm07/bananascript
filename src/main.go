package main

import (
	"bananascript/src/builtins"
	"bananascript/src/evaluator"
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"bananascript/src/repl"
	"fmt"
	"github.com/gookit/color"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		runFile(os.Args[1])
	} else {
		repl.Start()
	}
}

func runFile(fileName string) {
	theLexer, err := lexer.FromFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	theParser := parser.New(theLexer)
	context, environment := builtins.NewContextAndEnvironment()

	program, errors := theParser.ParseProgram(context)
	if len(errors) > 0 {
		errorStr := "Encountered %d error"
		if len(errors) > 1 {
			errorStr += "s"
		}
		errorStr += ":"
		fmt.Println(color.FgRed.Sprintf(errorStr, len(errors)))
		for _, err := range errors {
			fmt.Println(err.PrettyPrint(true))
		}
		os.Exit(1)
	} else {
		object := evaluator.Eval(program, environment)
		if err, isError := object.(*evaluator.ErrorObject); isError {
			fmt.Println(err.Message)
		}
	}
}
