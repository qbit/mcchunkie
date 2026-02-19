package chats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/mcstore"
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
	Type              string    `json:"type"`
	Short             bool      `json:"short"`
	Repo              string    `json:"repo"`
	AuthUser          string    `json:"auth_user"`
	AuthenticatedUser string    `json:"authenticated_user"`
	ID                string    `json:"id"`
	Author            Author    `json:"author,omitempty"`
	Committer         Committer `json:"committer"`
	Date              int       `json:"date"`
	ShortMessage      string    `json:"short_message"`
	Message           string    `json:"message"`
	Diffstat          Diffstat  `json:"diffstat,omitempty"`
}

func (n *Notification) String() string {
	// op committed got.git f9e653700..f9e653700^1 (main): fix gotd_parse_url() (https://git.gameoftrees.org/gitweb/?p=got.git;a=commitdiff;h=f9e653700)
	user := n.AuthenticatedUser
	if user == "" {
		user = n.AuthUser
	}
	return fmt.Sprintf("%s committed %s %s: %s (%s)",
		user,
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

func chatReplyv2(msgs GotNotifications, gotRoom string, cli Chat) error {
	str := []string{}

	for _, line := range msgs.Notifications {
		str = append(str, line.String())
	}

	return chatReply(strings.Join(str, "\n"), gotRoom, cli)
}

func chatReply(msg, gotRoom string, cli Chat) error {
	if gotRoom == "stdout" {
		log.Println(msg)
		return nil
	}
	for _, line := range strings.Split(msg, "\n") {
		log.Printf("GOT: sending '%s'\n", line)
		err := cli.Send(gotRoom, line)
		if err != nil {
			return fmt.Errorf("can not send commit info: %q", err)

		}
	}

	return nil
}

func GotListen(store *mcstore.MCStore, cli Chat) {
	var gotPort, err = store.Get("got_listen")
	if err != nil {
		log.Println(err)
		return
	}
	if gotPort != "" {
		var htpass, _ = store.Get("got_htpass")
		var gotRoom, _ = store.Get("got_room")

		log.Printf("GOT: listening on %q and sending messages to %q\n", gotPort, gotRoom)
		http.HandleFunc("/_got", func(w http.ResponseWriter, r *http.Request) {
			var msg string
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "got") {
				log.Printf("GOT: failed auth %q %q %q %q\n", user, pass, err, ok)
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

			err = chatReply(msg, gotRoom, cli)
			if err != nil {
				log.Printf("GOT: error sending: %q\n", err)
				http.Error(w, "unable to send", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "ok")

		})

		http.HandleFunc("/_got/v2", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				log.Printf("GOT: invalid method: '%q'\n", r.Method)
				http.Error(w, fmt.Sprintf("method %q not implemented", r.Method), http.StatusMethodNotAllowed)
				return
			}
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "got") {
				log.Printf("GOT: failed auth '%q %q'\n", user, pass)
				w.Header().Set("WWW-Authenticate", `Basic realm="got notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			gn := GotNotifications{}

			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("internal error: %s", err), http.StatusInternalServerError)
				return
			}

			_ = os.WriteFile("/tmp/mcchunkie-notification.json", data, 0600)

			err = json.Unmarshal(data, &gn)
			if err != nil {
				log.Printf("GOT: invalid data sent to server: '%s'\n", err)
				http.Error(w, fmt.Sprintf("invalid data sent to server: %s", err), http.StatusBadRequest)
				return
			}

			err = chatReplyv2(gn, gotRoom, cli)
			if err != nil {
				log.Printf("GOT: error sending commit info: '%s'\n", err)
				http.Error(
					w,
					fmt.Sprintf("can not send commit info: %s", err),
					http.StatusInternalServerError,
				)
				return
			}

		})

		log.Fatal(http.ListenAndServe(gotPort, nil))
	}
}
