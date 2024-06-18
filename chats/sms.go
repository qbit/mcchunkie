package chats

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
		smsPort, _      = store.Get("sms_listen")
		smsAllowed, _   = store.Get("sms_users")
		smtpUser, _     = store.Get("smtp_user")
		smtpReceiver, _ = store.Get("smtp_receiver")
		smsUsers        = strings.Split(smsAllowed, ",")
	)

	if smsPort != "" {
		var htpass, _ = store.Get("sms_htpass")

		log.Printf("SMS: listening on %q\n", smsPort)

		http.HandleFunc("/_sms", func(w http.ResponseWriter, r *http.Request) {
			var msg, from, id, date string
			emailSend := false
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
				id = r.URL.Query().Get("id")
				date = r.URL.Query().Get("date")
				msg = r.URL.Query().Get("message")
				from = r.URL.Query().Get("from")

				emailSend = true
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

						if emailSend {
							sc, err := smtp.Dial("localhost:25")
							if err != nil {
								log.Printf("SMS: smtp dial failed: %q\n", err)
								http.Error(w, "internal server error", http.StatusInternalServerError)
								return
							}

							if err := sc.Mail(smtpUser); err != nil {
								log.Printf("SMS: smtp Mail failed: %q\n", err)
								http.Error(w, "internal server error", http.StatusInternalServerError)
								return

							}
							if err := sc.Rcpt(smtpReceiver); err != nil {
								log.Printf("SMS: smtp Rcpt failed: %q\n", err)
								http.Error(w, "internal server error", http.StatusInternalServerError)
								return

							}

							wc, err := sc.Data()
							if err != nil {
								log.Printf("SMS: smtp Data failed: %q\n", err)
								http.Error(w, "internal server error", http.StatusInternalServerError)
								return

							}

							fmt.Fprintf(wc, fmt.Sprintf("To: %s\r\n", smtpReceiver))
							fmt.Fprintf(wc, fmt.Sprintf("From: %s\r\n", smtpUser))
							fmt.Fprintf(wc, fmt.Sprintf("Subject: Message received from number %s to number %s [%s]\r\n", from, from, id))
							fmt.Fprintf(wc, resp)

							defer wc.Close()

							err = sc.Quit()
							if err != nil {
								log.Printf("SMS: smtp Quit failed: %q\n", err)
								http.Error(w, "internal server error", http.StatusInternalServerError)
								return
							}
						} else {
							fmt.Fprint(w, resp)
						}
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
