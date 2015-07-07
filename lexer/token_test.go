package lexer

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}, args ...interface{}) {
	if actual != expected {
		msg := fmt.Sprint(args...)
		t.Errorf("Not equal(%s). actual: %#v expected: %#v", msg, actual, expected)
	}
}

var cases = map[string][]string{
	"(+ 12 (* 34 99))": []string{"(", "+", "12", "(", "*", "34", "99", ")", ")"},
	";":                []string{},
	"(display \"hello\")": []string{"(", "display", "\"hello\"", ")"},
}

func TestTokenizer(t *testing.T) {
	for input, tokens := range cases {
		tokenizer := NewTokenizer(strings.NewReader(input))

		for _, token := range tokens {
			tok, err := tokenizer.Next()
			assertEqual(t, err, nil)
			assertEqual(t, tok, token)
		}
		_, err := tokenizer.Next()
		assertEqual(t, err, io.EOF)
	}
}
