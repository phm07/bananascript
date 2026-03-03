//go:build js && wasm

package main

import (
	"bananascript/src/builtins"
	"bananascript/src/evaluator"
	"bananascript/src/lexer"
	"bananascript/src/parser"
	"strings"
	"syscall/js"
)

func main() {
	js.Global().Set("bananaRun", js.FuncOf(run))
	// Block forever to keep the WASM module alive
	select {}
}

func run(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return result("", nil)
	}

	code := args[0].String()

	var output strings.Builder
	printFn := func(s string) { output.WriteString(s) }
	promptFn := func(prompt string) string {
		return js.Global().Call("prompt", prompt).String()
	}

	context, environment := builtins.NewContextAndEnvironmentWithIO(printFn, promptFn)

	theLexer := lexer.FromCode(code)
	theParser := parser.New(theLexer)
	program, parseErrors := theParser.ParseProgram(context)

	if len(parseErrors) > 0 {
		errs := make([]interface{}, len(parseErrors))
		for i, e := range parseErrors {
			errs[i] = map[string]interface{}{
				"message": e.Message,
				"line":    e.Line,
				"col":     e.Col,
			}
		}
		return result(output.String(), errs)
	}

	obj := evaluator.Eval(program, environment)
	if errObj, isError := obj.(*evaluator.ErrorObject); isError {
		return result(output.String(), []interface{}{
			map[string]interface{}{
				"message": errObj.Message,
				"line":    0,
				"col":     0,
			},
		})
	}

	return result(output.String(), nil)
}

func result(output string, errors []interface{}) map[string]interface{} {
	if errors == nil {
		errors = []interface{}{}
	}
	return map[string]interface{}{
		"output": output,
		"errors": errors,
	}
}
