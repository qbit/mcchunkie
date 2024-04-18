package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestGotNotification(t *testing.T) {
	jsonStr := `{"notifications":[{"type":"commit", "short":false, "repo":"test-repo", "id":"34d7c970d4bd3a5832ecded86f2c28e83dbba2ba", "author":{ "full":"Flan Hacker <flan@openbsd.org>", "name":"Flan Hacker", "mail":"flan@openbsd.org", "user":"flan" }, "committer":{ "full":"Flan Hacker <flan@openbsd.org>", "name":"Flan Hacker", "mail":"flan@openbsd.org", "user":"flan" }, "date":"Thu Apr 18 16:47:40 2024 UTC", "short_message":"make changes", "message":"make changes\n", "diffstat":{ "files":[{ "action":"modified", "file":"alpha", "added":1, "removed":1 }], "total":{ "added":1, "removed":1 } }}]}`

	nots := &GotNotifications{}
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	err := dec.Decode(nots)

	if err != nil {
		t.Fatal(err)
	}

	for _, n := range nots.Notifications {
		fmt.Println(n.String())
	}
}
