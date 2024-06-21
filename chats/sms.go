package chats

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/plugins"
)

func smsCanSend(number string, numbers []string) bool {
	for _, s := range numbers {
		if number == s {
			return true
		}
	}
	return false
}

// SMSListen listens for our incoming sms
func SMSListen(store ChatStore, plugins *plugins.Plugins) {
	var (
		smsPort, _    = store.Get("sms_listen")
		smsAllowed, _ = store.Get("sms_users")
		smsUsers      = strings.Split(smsAllowed, ",")
	)

	if smsPort != "" {
		var htpass, _ = store.Get("sms_htpass")

		log.Printf("SMS: listening on %q\n", smsPort)

		http.HandleFunc("/_sms", func(w http.ResponseWriter, r *http.Request) {
			var msg, from string
			user, pass, ok := r.BasicAuth()
			if !ok {
				log.Println("SMS: basic auth no ok")
				w.Header().Set("WWW-Authenticate", `Basic realm="sms notify"`)
				http.Error(w, "auth error", http.StatusUnauthorized)
				return
			}

			if user != "sms" {
				log.Printf("SMS: failed auth for invalid user: %q, %q\n", user, pass)
				w.Header().Set("WWW-Authenticate", `Basic realm="sms notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if err != nil {
				log.Printf("SMS: failed auth %q %q\n", user, pass)
				w.Header().Set("WWW-Authenticate", `Basic realm="sms notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			switch r.Method {
			case http.MethodPost:
				err = r.ParseForm()
				if err != nil {
					http.Error(w, "invalid request", http.StatusBadRequest)
					return
				}
				msg = r.Form.Get("Body")
				from = r.Form.Get("From")
			case http.MethodGet:
				// to={TO}&from={FROM}&message={MESSAGE}&id={ID}&date={TIMESTAMP}
				msg = r.URL.Query().Get("message")
				from = r.URL.Query().Get("from")
			default:
				http.Error(
					w,
					fmt.Sprintf("method %q not implemented", r.Method),
					http.StatusMethodNotAllowed,
				)
				return
			}

			if smsCanSend(from, smsUsers) {
				msg = strings.TrimSuffix(msg, "\n")

				if msg == "" {
					fmt.Fprintf(w, "empty message")
					return
				}

				for _, p := range *plugins {
					if p.Match(from, msg) {
						log.Printf("%s: responding to '%s'", p.Name(), from)
						p.SetStore(store)

						resp := p.Process(from, msg)
						fmt.Fprint(w, resp)
					}
				}
			} else {
				log.Printf("number not allowed (%q)", from)
				http.Error(
					w,
					fmt.Sprintf("number not allowed (%q)", from),
					http.StatusMethodNotAllowed,
				)
				return
			}
		})
		log.Fatal(http.ListenAndServe(smsPort, nil))
	}
}
