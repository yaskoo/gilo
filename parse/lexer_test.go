package parse

import (
	"testing"
)

func TestNext(t *testing.T) {
	if lex("").next() != eof {
		t.Error("should return eof rune on empty input")
	}

	l := lex("golang")
	r := l.next()
	if r != 'g' {
		t.Error("should return the next rune")
	}
	if l.position <= 0 || l.width <= 0 || l.position != (l.start+l.width) {
		t.Errorf("should advance position and update width: start: %d, position: %d, width: %d\n", l.start, l.position, l.width)
	}

	l.next() // o
	l.next() // l
	l.next() // a
	l.next() // n
	l.next() // g

	r = l.next()
	if r != eof {
		t.Error("should return eof when end of input is reached")
	}

	if l.width != 0 {
		t.Error("width should become zero when we reach end of input string")
	}
}

func TestBackup(t *testing.T) {
	l := lex("golang")
	l.backup()
	if l.position < 0 {
		t.Error("position shold not become negative")
	}

	l.next()
	l.backup()
	if l.width != 0 {
		t.Error("width should become zero after backup")
	}
	if l.next() != 'g' {
		t.Error("previous rune should be unread")
	}

	l = lex("golang")
	l.start = 1
	l.position = 1
	l.backup()
	if l.position < l.start {
		t.Error("position should not be less than start")
	}
}

func TestPeek(t *testing.T) {
	l := lex("golang")
	r := l.peek()
	if r != 'g' {
		t.Error("should return the next rune")
	}

	if l.position != 0 {
		t.Error("peek should not advance position")
	}

	if l.width != 0 {
		t.Error("peek should not remember next rune width")
	}
}

func TestEmit(t *testing.T) {
	l := lex("golang")

	go func() {
		l.next()
		l.next()

		l.emit(textType)
	}()

	tok := <-l.tokens
	if tok.typ != textType {
		t.Error("expecting text token")
	}
	if tok.val != "go" {
		t.Error("expecting 'go' value")
	}
	if l.start != l.position {
		t.Error("start and position should be equal after emit")
	}
}

func TestIgnore(t *testing.T) {
	l := lex("golang")
	l.next()
	l.next()
	l.ignore()

	if l.start != 2 || l.position != 2 {
		t.Error("start and position should be equal after ignore")
	}

	if r := l.next(); r != 'l' {
		t.Errorf("expected 'l', got: %s", string(r))
	}

}
