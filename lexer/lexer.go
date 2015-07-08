package lexer

import (
	"fmt"
	"io"

	"github.com/Xuyuanp/common"
)

type stateFunc func(*Lexer) stateFunc

type Lexer struct {
	name       string
	tokenizer  *Tokenizer
	state      stateFunc
	output     chan interface{}
	mainStack  common.Stack
	stateStack common.Stack
}

func NewLexer(name string, r io.Reader) *Lexer {
	return &Lexer{
		name:       name,
		tokenizer:  NewTokenizer(r),
		output:     make(chan interface{}, 10),
		state:      idleState,
		mainStack:  common.NewStack(),
		stateStack: common.NewStack(),
	}
}

func (l *Lexer) Output() <-chan interface{} {
	return l.output
}

func (l *Lexer) Run() {
	for ; l.state != nil; l.state = l.state(l) {
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFunc {
	return func(lexer *Lexer) stateFunc {
		msg := fmt.Sprintf(format, args...)
		msg = fmt.Sprintf("[Lexer: %s] %s", lexer.name, msg)
		l.output <- msg
		return nil
	}
}

func idleState(l *Lexer) stateFunc {
	token, err := l.tokenizer.Next()
	if err != nil {
		return l.errorf(err.Error())
	}
	if token.typ == typeInteger || token.typ == typeFloat {
		l.output <- token.val
		return idleState
	}

	if token.typ == typeLeftParen {
		l.mainStack.Push(token)
		return inPrecedureState
	}

	return l.errorf("invlid token %s", token)
}

func inPrecedureState(l *Lexer) stateFunc {
	token, err := l.tokenizer.Next()
	for ; err == nil; token, err = l.tokenizer.Next() {
		switch token.typ {
		case typeRightParen:
			return endPrecedureState
		case typeLeftParen:
			l.mainStack.Push(token)
			l.stateStack.Push(stateFunc(inPrecedureState))
			return inPrecedureState
		default:
			l.mainStack.Push(token)
		}
	}

	if err != nil {
		return l.errorf(err.Error())
	}

	return idleState
}

func endPrecedureState(l *Lexer) stateFunc {
	args := []*Token{}
	for l.mainStack.Len() > 0 {
		token := l.mainStack.Pop().(*Token)
		if token.typ == typeLeftParen {
			if len(args) == 0 {
				return l.errorf("empty ()")
			}
			result, err := doPrecedure(args...)
			if err != nil {
				return l.errorf(err.Error())
			}

			if s := l.stateStack.Pop(); s != nil {
				l.mainStack.Push(result)
				return s.(stateFunc)
			}
			l.output <- result.val
			return idleState
		}
		args = append([]*Token{token}, args...)
	}
	return l.errorf("missing (")
}
