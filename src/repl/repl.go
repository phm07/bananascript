package repl

import (
	"bananascript/src/builtins"
	"bananascript/src/evaluator"
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"bufio"
	"fmt"
	"os"
)

const PROMPT = "> "

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	context, environment := builtins.NewContextAndEnvironment()

	for {
		fmt.Printf(PROMPT)

		if !scanner.Scan() {
			return
		}

		input := scanner.Text() + ";"

		theLexer := lexer.FromCode(input)
		theParser := parser.New(theLexer)

		program, errors := theParser.ParseProgram(context)
		newContext := program.Context
		newEnvironment := evaluator.ExtendEnvironment(environment, newContext)

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err.PrettyPrint(false))
			}
		} else {
			var result evaluator.Object = nil
			for _, statement := range program.Statements {
				result = evaluator.Eval(statement, newEnvironment)
			}
			if result != nil {
				fmt.Println(result.ToString())
			}
			if _, isError := result.(*evaluator.ErrorObject); !isError {
				context = newContext
				environment = newEnvironment
			}
		}
	}
}
