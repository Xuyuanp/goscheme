package lexer

import (
	"bufio"
	"bytes"
	"io"
)

type Tokenizer struct {
	tokens chan string
	err    chan error
	buffer *bytes.Buffer
	reader *bufio.Reader
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		tokens: make(chan string, 3),
		err:    make(chan error, 2),
		buffer: bytes.NewBuffer([]byte{}),
		reader: bufio.NewReader(r),
	}
}

func (t *Tokenizer) Next() (string, error) {
	for {
		select {
		case token := <-t.tokens:
			return token, nil
		case err := <-t.err:
			return "", err
		default:
			t.scanStream()
		}
	}
}

func (t *Tokenizer) emit() {
	if t.buffer.Len() > 0 {
		t.tokens <- t.buffer.String()
		t.buffer.Reset()
	}
}

func (t *Tokenizer) scanStream() {
	c, err := t.reader.ReadByte()
	if err != nil {
		t.err <- err
		return
	}
	switch c {
	case ';':
	case ' ', '\t', '\n', '\r':
		t.emit()
	case '(', ')':
		t.emit()
		t.buffer.WriteByte(c)
		t.emit()
	case '"':
		t.emit()
		t.buffer.WriteByte(c)
		str, err := t.reader.ReadBytes('"')
		if err != nil {
			t.err <- err
			return
		}
		t.buffer.Write(str)
		t.emit()
	default:
		t.buffer.WriteByte(c)
	}
}
