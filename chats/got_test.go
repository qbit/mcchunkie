package chats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGotNotification(t *testing.T) {
	jsonData, err := os.ReadFile("test_body.json")
	if err != nil {
		t.Fatal(err)
	}

	nots := &GotNotifications{}
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	err = dec.Decode(nots)
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range nots.Notifications {
		fmt.Println(n.String())
	}
}
