package main

import (
	"fmt"
	"os"

	"github.com/Youssef-Mak/baby-interpreter/pkg/repl"
)

const BABY = `
(  _ \  /__\  (  _ \( \/ )  
 ) _ < /(__)\  ) _ < \  /   
(____/(__)(__)(____/ (__)   
`

func main() {

	fmt.Println("Baby Version 1.0.0")
	fmt.Println(BABY)

	fmt.Println("\nTo import a baby file(*.bb) simply input the filename with the .bb extension. (Ex: >> filename.bb)")

	repl.Initialize(os.Stdin, os.Stdout)

}
