package plugins

import (
	"testing"
)

func TestBeatMatch(t *testing.T) {
	testStrings := make(map[string]bool)
	testStrings["what time is it??!!?!"] = true
	testStrings["what time is it"] = false
	testStrings["man: 2 pledge"] = false
	testStrings[".beat"] = true
	testStrings["beattime?"] = true
	testStrings["beat time?"] = true

	b := &Beat{}
	for msg, should := range testStrings {
		if b.Match("", msg) != should {
			t.Errorf("Beat expected to match %q (%t); but doesn't\n", msg, should)
		}
	}
}
