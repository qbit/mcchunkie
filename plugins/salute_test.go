package plugins

import (
	"testing"
)

func TestSaluteMatch(t *testing.T) {
	testStrings := make(map[string]bool)
	testStrings["o7"] = true
	testStrings[`7o`] = false
	testStrings["_o"] = false
	testStrings["o_"] = false

	h := &Salute{}
	for msg, should := range testStrings {
		if h.Match("", msg) != should {
			t.Errorf("Salute expected to match %q (%t); but doesn't\n", msg, should)
		}
	}
}
