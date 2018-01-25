package parse

import (
	"fmt"
)

// Message contains a parsed message
type Message struct {
	Subject string
	Body    string
	Tags    []string
}

func (m Message) String() string {
	return fmt.Sprintf("Subject: %s\nBody: %s\nTags: %v", m.Subject, m.Body, m.Tags)
}

// Parser is responsible for parsing a commit message
// It extracts the subject line, message body, tags, etc
type Parser struct {
	lex *lexer
	Msg Message
}

// Parse starts the parsing process
func (p *Parser) Parse() {
	go p.lex.run()

	var body bool
	for tok := range p.lex.tokens {
		switch tok.typ {
		case tagType:
			p.Msg.Tags = append(p.Msg.Tags, tok.val)
			fallthrough
		case textType:
			if body {
				p.Msg.Body += tok.val
			} else {
				p.Msg.Subject += tok.val
			}
		case newLineType:
			body = true
		}
	}

	fmt.Println(p.Msg)
}

// NewParser contructs a new parsing instance initialized w/ the provided input
// Use Parse to start parsing
func NewParser(input string) *Parser {
	return &Parser{
		lex: lex(input),
	}
}
