package main

import (
	"testing"
)

func TestParseErrata(t *testing.T) {
	got, err := ParseRemoteErrata("http://ftp.openbsd.org/pub/OpenBSD/patches/6.6/common/")
	if err != nil {
		t.Error(err)
	}
	l := len(got.List)
	if l < 1 {
		t.Errorf("errata count %d; want > 1", l)
	}

	erratum := got.List[len(got.List)-1]

	err = erratum.Fetch()
	if err != nil {
		t.Errorf("can't fetch data for erratum\n%s", err)
	}
}
