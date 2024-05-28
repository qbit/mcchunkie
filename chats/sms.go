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
	var smsPort, _ = store.Get("sms_listen")
	var smsAllowed, _ = store.Get("sms_users")
	var smsUsers = strings.Split(smsAllowed, ",")

	if smsPort != "" {
		var htpass, _ = store.Get("sms_htpass")

		log.Printf("SMS: listening on %q\n", smsPort)

		http.HandleFunc("/_sms", func(w http.ResponseWriter, r *http.Request) {
			var msg, from string
			if r.Method != http.MethodPost {
				log.Printf("SMS: invalid method: '%q'\n", r.Method)
				http.Error(w, fmt.Sprintf("method %q not implemented", r.Method), http.StatusMethodNotAllowed)
				return
			}
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

			err = r.ParseForm()
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}

			log.Println(r.Method)

			switch r.Method {
			case http.MethodPost:
				msg = r.Form.Get("Body")
				from = r.Form.Get("From")
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
