package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/matrix-org/gomatrix"
	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/plugins"
)

type Notification struct {
	Short    bool   `json:"short"`
	ID       string `json:"id"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Message  string `json:"message"`
	Diffstat struct {
	} `json:"diffstat"`
	Changes struct {
	} `json:"changes"`
}

func (n *Notification) String() string {
	// op committed got.git f9e653700..f9e653700^1 (main): fix gotd_parse_url() (https://git.gameoftrees.org/gitweb/?p=got.git;a=commitdiff;h=f9e653700)
	return fmt.Sprintf("%s committed %s: %s (%s)",
		n.Author,
		n.ID,
		n.Message,
		fmt.Sprintf("https://git.gameoftrees.org/gitweb/?p=%s;a=commitdiff;h=%s",
			"repo",
			n.ID),
	)
}

type GotNotifications struct {
	Notifications []Notification `json:"notifications"`
}

func gotListen(store *FStore, cli *gomatrix.Client) {
	var gotPort, _ = store.Get("got_listen")
	if gotPort != "" {
		var htpass, _ = store.Get("got_htpass")
		var gotRoom, _ = store.Get("got_room")

		log.Printf("GOT: listening on %q and sending messages to %q\n", gotPort, gotRoom)

		http.HandleFunc("/_got/v2", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, fmt.Sprintf("method %q not implemented", r.Method), http.StatusMethodNotAllowed)
				return
			}
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "got") {
				log.Printf("GOT: failed auth '%s'\n", user)
				w.Header().Set("WWW-Authenticate", `Basic realm="got notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			gn := GotNotifications{}

			dec := json.NewDecoder(r.Body)
			err = dec.Decode(&gn)
			if err != nil {
				http.Error(w, "invalid data sent to server", http.StatusBadRequest)
				return
			}
			for _, line := range gn.Notifications {
				log.Printf("GOT: sending '%s'\n", line.String())
				err = plugins.SendUnescNotice(cli, gotRoom, line.String())
				if err != nil {
					http.Error(
						w,
						fmt.Sprintf("can not send commit info: %s", err),
						http.StatusInternalServerError,
					)
					return
				}
			}
		})

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
