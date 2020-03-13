package plugins

import (
	"testing"
)

func TestHighFiveMatch(t *testing.T) {
	testStrings := make(map[string]bool)
	testStrings["o/"] = true
	testStrings[`\o`] = true
	testStrings["_o"] = false
	testStrings["o_"] = false

	h := &HighFive{}
	for msg, should := range testStrings {
		if h.Match("", msg) != should {
			t.Errorf("HighFive expected to match %q (%t); but doesn't\n", msg, should)
		}
	}
}
