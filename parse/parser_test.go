package parse

import (
	"testing"
)

func TestParser(t *testing.T) {
	p := NewParser("go is #love, go is #life\n this is #ze#body")

	p.Parse()
}
