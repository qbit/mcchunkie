package plugins

import (
	"testing"
)

func TestPluginsToMe(t *testing.T) {
	if ToMe("a", "a") == false {
		t.Error("ToMe expected true; got false\n")
	}
}

func TestPluginsNameRE(t *testing.T) {
	n := NameRE.ReplaceAllString("@test:test.com", "$1")
	if n != "test" {
		t.Errorf("NameRE expected 'test'; got %q\n", n)
	}
}
