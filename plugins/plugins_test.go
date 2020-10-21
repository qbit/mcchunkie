package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
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

func TestPluginsRemoveName(t *testing.T) {
	expected := "this is for you"
	n := RemoveName("mctest", fmt.Sprintf("mctest: %s", expected))
	if n != expected {
		t.Errorf("Expected %q; got %q\n", expected, n)
	}
}

type testResp struct {
	Name string `json:"test"`
}

func TestHTTPRequestDoJSON(t *testing.T) {
	if runtime.GOOS != "plan9" {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprintln(w, `{"test":"success"}`)
			if err != nil {
				t.Error(err)
			}
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
}
