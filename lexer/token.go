package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type tokenType int

func (tt tokenType) String() string {
	switch tt {
	case typeUnknown:
		return "unknown"
	case typeLeftParen:
		return "left paren"
	case typeRightParen:
		return "right paren"
	case typeInteger:
		return "integer"
	case typeFloat:
		return "float"
	case typeString:
		return "string"
	case typePrecedure:
		return "precedure"
	case typeDefine:
		return "define"
	}
	return fmt.Sprintf("[token type %d]", tt)
}

const (
	typeUnknown tokenType = iota
	typeLeftParen
	typeRightParen
	typeInteger
	typeFloat
	typeString
	typePrecedure
	typeDefine
)

const (
	leftParen  = '('
	rightParen = ')'
	comment    = ';'
	quote      = '"'
)

type Token struct {
	typ tokenType
	val interface{}
	raw string
}

type Tokenizer struct {
	tokens chan *Token
	err    chan error
	buffer *bytes.Buffer
	reader *bufio.Reader
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		tokens: make(chan *Token, 3),
		err:    make(chan error, 2),
		buffer: bytes.NewBuffer([]byte{}),
		reader: bufio.NewReader(r),
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	for {
		select {
		case token := <-t.tokens:
			return token, nil
		case err := <-t.err:
			return nil, err
		default:
			t.scanStream()
		}
	}
}

func (t *Tokenizer) emit() {
	if t.buffer.Len() > 0 {
		token := t.parseToken(t.buffer.String())
		t.tokens <- token
		t.buffer.Reset()
	}
}

func (t *Tokenizer) parseToken(raw string) *Token {
	var typ = typeUnknown
	var val interface{}
	switch {
	case raw == "define":
		typ = typeDefine
	case raw == "(":
		typ = typeLeftParen
		val = raw
	case raw == ")":
		typ = typeRightParen
		val = raw
	case isString(raw):
		typ = typeString
		val = raw[1 : len(raw)-1]
	case isFloat(raw):
		val, _ = strconv.ParseFloat(raw, 64)
		typ = typeFloat
	case isInteger(raw):
		val, _ = strconv.Atoi(raw)
		typ = typeInteger
	default:
		var ok bool
		if val, ok = precedureMap[raw]; ok {
			typ = typePrecedure
		}

	}
	token := &Token{
		typ: typ,
		raw: raw,
		val: val,
	}
	return token
}

func (t *Tokenizer) scanStream() {
	c, err := t.reader.ReadByte()
	if err != nil {
		t.err <- err
		return
	}
	switch c {
	case comment:
	case ' ', '\t', '\n', '\r':
		t.emit()
	case leftParen, rightParen:
		t.emit()
		t.buffer.WriteByte(c)
		t.emit()
	case quote:
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
