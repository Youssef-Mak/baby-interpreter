package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/Youssef-Mak/baby-interpreter/pkg/evaluator"
	"github.com/Youssef-Mak/baby-interpreter/pkg/object"
	"github.com/Youssef-Mak/baby-interpreter/pkg/parser"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
)

const PROMPT = ">> "

func Initialize(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		evaluated := object.Object(nil)
		ok := false

		if !scanned {
			return
		}

		line := scanner.Text()
		if isBabyFile(line) {
			buf, err := ioutil.ReadFile(line)
			if err == nil {
				contents := string(buf)
				io.WriteString(out, contents)
				io.WriteString(out, "\n")
				evaluated, ok = InterpretInput(contents, out, env)
			} else {
				io.WriteString(out, fmt.Sprintf("Error reading Baby File: %s", err.Error()))
				io.WriteString(out, "\n")
				continue
			}
		} else {
			evaluated, ok = InterpretInput(line, out, env)
		}

		if !ok {
			continue
		}
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func InterpretInput(input string, out io.Writer, env *object.Environment) (object.Object, bool) {
	tokenizer := tokenizer.New(input)
	parser := parser.New(tokenizer)

	program := parser.ParseProgram()
	if len(parser.GetErrors()) != 0 {
		printParserErrors(out, parser.GetErrors())
		return nil, false
	}

	evaluated := evaluator.Eval(program, env)
	return evaluated, true
}

func isBabyFile(toMatch string) bool {
	match, _ := regexp.Match(`[^\\]*\.(bb)`, []byte(toMatch))
	return match
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
