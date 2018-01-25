package parse

import (
	"fmt"
	"unicode/utf8"
)

const eof = -1

// lexer is responsible for extracting tokens from an input string
type lexer struct {
	start    int
	position int
	width    int
	input    string
	tokens   chan token
}

// custom token type which is just an int w/ some helper functions
type tokenType int

func (t tokenType) String() string {
	switch t {
	case textType:
		return fmt.Sprintf("TEXT(%d)", t)
	case tagType:
		return fmt.Sprintf("TAG(%d)", t)
	case newLineType:
		return fmt.Sprintf("NL(%d)", t)
	case eofType:
		return fmt.Sprintf("EOF(%d)", t)
	case errType:
		return fmt.Sprintf("ERROR(%d)", t)
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

// token types
const (
	_ tokenType = iota
	textType
	tagType
	newLineType
	eofType
	errType
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch {
	case t.typ == eofType:
		return "EOF"
	case t.typ == errType:
		return t.val
	case t.typ == newLineType:
		return "NL"
	case len(t.val) > 15:
		return fmt.Sprintf("%.15q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

// emit sends a token to the tokens channel
// if start and position are equal nothing is sent
func (l *lexer) emit(typ tokenType) {
	l.tokens <- token{
		typ: typ,
		val: l.input[l.start:l.position],
	}

	l.start = l.position
}

// next returns the next rune in the input string
// return eof if the postion becomes >= to the length of the input string
func (l *lexer) next() rune {
	if l.position >= len(l.input) {
		l.width = 0
		return eof
	}

	r, s := utf8.DecodeRuneInString(l.input[l.position:])
	l.width = s
	l.position += s
	return r
}

// backup unreads the last read rune
// successive calls have no afect
func (l *lexer) backup() {
	l.position -= l.width
	l.width = 0
	if l.position < l.start {
		l.position = l.start
	}
}

// peek returns the next rune in the input string w/o advancing the position
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// ignore discards what has been read so far
func (l *lexer) ignore() {
	l.start = l.position
}

// run starts the lexing process looping thu state functions until one retuns nil
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// lex initializes a new lexer using the provided input string
func lex(input string) *lexer {
	return &lexer{
		input:  input,
		tokens: make(chan token),
	}
}
