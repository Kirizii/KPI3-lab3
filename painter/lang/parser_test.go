package lang

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	script := `
white
bgrect 0.1 0.1 0.9 0.9
figure 0.5 0.5
green
move 0.6 0.6
update
reset
`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(script))
	if err != nil {
		t.Fatal(err)
	}

	if len(ops) != 7 {
		t.Errorf("Expected 7 operations, got %d", len(ops))
	}
}
