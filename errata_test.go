package main

import (
	"testing"
)

func TestParseErrata(t *testing.T) {
	got, err := ParseRemoteErrata("https://www.openbsd.org/errata66.html")
	if err != nil {
		t.Error(err)
	}
	l := len(got.List)
	if l == 0 {
		t.Errorf("errata count %d; want > 0", l)
	}
}
