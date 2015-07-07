package lexer

import (
	"fmt"
	"io"
	"log"
)

type stateFunc func(*Lexer) stateFunc

type Lexer struct {
	name      string
	tokenizer *Tokenizer
	state     stateFunc
}

func NewLexer(name string, r io.Reader) *Lexer {
	return &Lexer{
		name:      name,
		tokenizer: NewTokenizer(r),
	}
}

func (l *Lexer) run() {
	for ; l.state != nil; l.state = l.state(l) {
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFunc {
	return func(lexer *Lexer) stateFunc {
		msg := fmt.Sprintf(format, args...)
		log.Printf("[Lexer: %s]%s", lexer.name, msg)
		return nil
	}
}
