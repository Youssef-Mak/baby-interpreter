package repl

import (
	"bufio"
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/evaluator"
	"github.com/Youssef-Mak/baby-interpreter/pkg/object"
	"github.com/Youssef-Mak/baby-interpreter/pkg/parser"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
	"io"
)

const PROMPT = ">> "

func Initialize(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		tokenizer := tokenizer.New(line)
		parser := parser.New(tokenizer)

		program := parser.ParseProgram()
		if len(parser.GetErrors()) != 0 {
			printParserErrors(out, parser.GetErrors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
