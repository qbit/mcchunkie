package plugins

import (
	"testing"
)

func TestFederRE(t *testing.T) {
	testStrings := make(map[string]bool)
	testStrings["feder: what.com"] = true
	testStrings["tayshame: what.com"] = true
	testStrings["feder:what"] = false
	testStrings["tayshame:what"] = false

	h := &Feder{}
	for msg, should := range testStrings {
		if h.Match("", msg) != should {
			t.Errorf("HighFive expected to match %q (%t); but doesn't\n", msg, should)
		}
	}
}
