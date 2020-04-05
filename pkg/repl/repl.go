package repl

import (
	"bufio"
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
	"io"
)

const PROMPT = ">> "

func Initialize(in io.Reader) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		tokenizer := tokenizer.New(line)

		for tok := tokenizer.NextToken(); tok.Type != token.EOF; tok = tokenizer.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
