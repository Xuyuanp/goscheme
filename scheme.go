package main

import (
	"os"

	"github.com/Xuyuanp/goscheme/repl"
)

func main() {
	if len(os.Args) < 2 {
		// TODO: REPL
		repl.Start()
		return
	}

}
