package plugins

import (
	"fmt"
	"testing"
)

func TestOpenBSDManfix(t *testing.T) {
	testStrings := make(map[string]string)
	testStrings["man: pledge"] = "https://man.openbsd.org/pledge"
	testStrings["man: 2 pledge"] = "https://man.openbsd.org/pledge.2"
	testStrings["man: unveil"] = "https://man.openbsd.org/unveil"
	testStrings["man: 3p vars"] = "https://man.openbsd.org/man3p/vars.3p"

	om := &OpenBSDMan{}
	for msg, resp := range testStrings {
		matched := fmt.Sprintf("https://man.openbsd.org/%s", om.fix(msg))
		if matched != resp {
			t.Errorf("OpenBSDMan expected %q; got %q (%q)", resp, matched, msg)
		}
	}
}
