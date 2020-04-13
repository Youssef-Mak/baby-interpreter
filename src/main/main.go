package main

import (
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/repl"
	"os"
)

const BABY = `
(  _ \  /__\  (  _ \( \/ )  
 ) _ < /(__)\  ) _ < \  /   
(____/(__)(__)(____/ (__)   
`

func main() {

	fmt.Println("Baby Version 0.0.0")
	fmt.Println(BABY)

	repl.Initialize(os.Stdin, os.Stdout)

}
