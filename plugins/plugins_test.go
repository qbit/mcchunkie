package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

type testResp struct {
	Name string `json:"test"`
}

func TestHTTPRequestDoJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"test":"success"}`)
	}))
	defer ts.Close()

	var tr = &testResp{}

	req := HTTPRequest{
		ResBody: tr,
		URL:     ts.URL,
	}

	err := req.DoJSON()
	if err != nil {
		t.Error(err)
	}

	if tr.Name != "success" {
		t.Errorf("Expected 'test'; got '%s'\n", tr.Name)
	}
}
