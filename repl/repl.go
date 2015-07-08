package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Xuyuanp/goscheme/lexer"
)

var lineNumber = 0

const promptFormat = "[%d]==> "

func Start() error {
	repl := NewRepl()
	return repl.Run()
}

type Repl struct {
	w          io.Writer
	lineNumber int
	lexer      *lexer.Lexer
}

func NewRepl() *Repl {
	r, w := io.Pipe()
	return &Repl{
		w:     w,
		lexer: lexer.NewLexer("REPL", r),
	}
}

func (repl *Repl) Run() error {
	go repl.lexer.Run()
	for {
		select {
		case output := <-repl.lexer.Output():
			repl.print(output)
		default:
			repl.waitInput()
		}
	}
	return nil
}

func (repl *Repl) print(output interface{}) {
	fmt.Println(output)
	repl.lineNumber++
}

func (repl *Repl) waitInput() {
	line := repl.getLine()
	fmt.Fprintln(repl.w, line)
}

var scanner = bufio.NewScanner(os.Stdin)

func (repl *Repl) getLine() (line string) {
	fmt.Printf(promptFormat, lineNumber)
	scanner.Scan()
	line = scanner.Text()
	return
}
