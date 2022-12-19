package repl

import (
	"bananascript/src/builtins"
	"bananascript/src/evaluator"
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"bananascript/src/types"
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

		newContext := types.ExtendContext(context)
		program, errors := theParser.ParseProgram(newContext)

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err.PrettyPrint(false))
			}
		} else {
			context = newContext
			environment = evaluator.ExtendEnvironment(environment, types.NewContext())
			result := evaluator.Eval(program, environment)
			if result != nil {
				fmt.Println(result.ToString())
			}
		}
	}
}
