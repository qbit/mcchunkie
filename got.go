package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/matrix-org/gomatrix"
	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/plugins"
)

func gotListen(store *FStore, cli *gomatrix.Client) {
	var gotPort, _ = store.Get("got_listen")
	if gotPort != "" {
		var htpass, _ = store.Get("got_htpass")
		var gotRoom, _ = store.Get("got_room")

		log.Printf("GOT: listening on %q and sending messages to %q\n", gotPort, gotRoom)

		http.HandleFunc("/_got", func(w http.ResponseWriter, r *http.Request) {
			var msg string
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "got") {
				log.Printf("GOT: failed auth '%s'\n", user)
				w.Header().Set("WWW-Authenticate", `Basic realm="got notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = r.ParseForm()
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}

			switch r.Method {
			case http.MethodGet:
				msg = r.Form.Get("message")
			case http.MethodPost:
				msg = r.Form.Get("file")
			default:
				http.Error(
					w,
					fmt.Sprintf("method %q not implemented", r.Method),
					http.StatusMethodNotAllowed,
				)
				return
			}

			msg = strings.TrimSuffix(msg, "\n")

			if msg == "" {
				fmt.Fprintf(w, "empty message")
				return
			}

			for _, line := range strings.Split(msg, "\n") {
				log.Printf("GOT: sending '%s'\n", line)
				err = plugins.SendUnescNotice(cli, gotRoom, line)
				if err != nil {
					http.Error(
						w,
						fmt.Sprintf("can not send commit info: %s", err),
						http.StatusInternalServerError,
					)
					return
				}
			}

			fmt.Fprintf(w, "ok")

		})

		log.Fatal(http.ListenAndServe(gotPort, nil))
	}
}
