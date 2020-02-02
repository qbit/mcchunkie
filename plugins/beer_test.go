package plugins

import (
	"fmt"
	"testing"
)

func TestBeer(t *testing.T) {
	beer := &Beer{}
	b, err := beer.get("oskar blues")
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	if b.Nhits == 0 {
		t.Errorf("Expected 7 results; got %d\n", b.Nhits)
	}
	fmt.Printf("%+v", b)
}
