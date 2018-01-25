package parse

import (
	"testing"
)

func TestLexText(t *testing.T) {
	l := lex("golang\n is awesome")
	go l.run()

	tok := <-l.tokens
	if tok.typ != textType {
		t.Errorf("expected token type text, got: %s", tok.typ)
	}

	if tok.val != "golang" {
		t.Errorf("exepected 'golang' token value, got: %s", tok.val)
	}

	tok = <-l.tokens
	if tok.typ != newLineType {
		t.Errorf("expected token type nl, got: %s", tok.typ)
	}

	tok = <-l.tokens
	if tok.typ != textType {
		t.Errorf("expected token type text, got: %s", tok.typ)
	}

	if tok.val != " is awesome" {
		t.Errorf("exepected ' is awesome' token value, got: %s", tok.val)
	}

	tok = <-l.tokens
	if tok.typ != eofType {
		t.Errorf("expected token type eof, got: %s", tok.typ)
	}
}

func TestLexTag(t *testing.T) {
	l := lex("golang is #awesome")
	go l.run()

	tok := <-l.tokens
	if tok.typ != textType {
		t.Errorf("expected token type text, got: %s", tok.typ)
	}

	if tok.val != "golang is " {
		t.Errorf("exepected 'golang is ' token value, got: %s", tok.val)
	}

	tok = <-l.tokens
	if tok.typ != tagType {
		t.Errorf("expected token type tag, got: %s", tok.typ)
	}

	if tok.val != "#awesome" {
		t.Errorf("exepected '#awesome' token value, got: %s", tok.val)
	}

	tok = <-l.tokens
	if tok.typ != eofType {
		t.Errorf("expected token type eof, got: %s", tok.typ)
	}
}
