package main

import (
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/repl"
	"os"
)

func main() {

	fmt.Println("Baby Version 0.0.0")

	repl.Initialize(os.Stdin)

}
