package main

import (
	"bananascript/src/builtins"
	"bananascript/src/evaluator"
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"bananascript/src/repl"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
)

func main() {
	help := flag.Bool("help", false, "show help")
	forceColor := flag.Bool("forceColor", false, "force colorized output")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *forceColor {
		color.ForceColor()
	}

	if flag.NArg() > 0 {
		runFile(flag.Arg(0))
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
