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

type Author struct {
	Full string `json:"full"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	User string `json:"user"`
}
type Committer struct {
	Full string `json:"full"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	User string `json:"user"`
}
type Files struct {
	Action  string `json:"action"`
	File    string `json:"file"`
	Added   int    `json:"added"`
	Removed int    `json:"removed"`
}
type Total struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
}
type Diffstat struct {
	Files []Files `json:"files"`
	Total Total   `json:"total"`
}
type Notification struct {
	Type         string    `json:"type"`
	Short        bool      `json:"short"`
	Repo         string    `json:"repo"`
	ID           string    `json:"id"`
	Author       Author    `json:"author,omitempty"`
	Committer    Committer `json:"committer"`
	Date         string    `json:"date"`
	ShortMessage string    `json:"short_message"`
	Message      string    `json:"message"`
	Diffstat     Diffstat  `json:"diffstat,omitempty"`
}

func (n *Notification) String() string {
	// op committed got.git f9e653700..f9e653700^1 (main): fix gotd_parse_url() (https://git.gameoftrees.org/gitweb/?p=got.git;a=commitdiff;h=f9e653700)
	return fmt.Sprintf("%s committed %s %s: %s (%s)",
		n.Committer.User,
		n.Repo,
		n.ID,
		n.ShortMessage,
		fmt.Sprintf("https://got.gameoftrees.org/?action=diff&commit=%s&headref=HEAD&path=%s",
			n.ID,
			n.Repo,
		),
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
				log.Printf("GOT: invalid data sent to server: '%s'\n", err)
				http.Error(w, fmt.Sprintf("invalid data sent to server: %s", err), http.StatusBadRequest)
				return
			}
			for _, line := range gn.Notifications {
				log.Printf("GOT: sending '%s'\n", line.String())
				err = plugins.SendUnescNotice(cli, gotRoom, line.String())
				if err != nil {
					log.Printf("GOT: error sending commit info: '%s'\n", err)
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
