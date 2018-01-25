package parse

import (
	_ "strings"
	"unicode"
)

const whitespace = "\t\r\n "

type stateFn func(*lexer) stateFn

func lexText(l *lexer) stateFn {
	// TODO: add escape mechanism \# \[ etc
	for {
		r := l.next()

		if r == '#' {
			l.backup()
			if l.start < l.position {
				l.emit(textType)
			}
			return lexTag
		}

		if r == '\n' {
			l.backup()
			l.emit(textType)
			l.next()
			l.emit(newLineType)
			continue
		}

		if r == eof {
			if l.start < l.position {
				l.emit(textType)
			}
			break
		}
	}

	if l.start == l.position {
		l.emit(eofType)
	} else {
		// TODO: error
	}
	return nil
}

func lexTag(l *lexer) stateFn {
	l.next()
	for {
		r := l.next()
		if r == eof {
			// a.k.a. we've lexed more than the hashtag (#)
			if l.position-l.start > 1 {
				l.emit(tagType)
			} else {
				// TODO: error
			}
			break
		}

		if r != '_' && r != '-' && (unicode.IsPunct(r) || unicode.IsSpace(r)) {
			l.backup()
			l.emit(tagType)
			break
		}
	}
	return lexText
}
