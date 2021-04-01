package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/plugins"
)

func smsListen(store *FStore, plugins *plugins.Plugins) {
	var smsPort, _ = store.Get("sms_listen")
	if smsPort != "" {
		var htpass, _ = store.Get("sms_htpass")

		log.Printf("SMS: listening on %q\n", smsPort)

		http.HandleFunc("/_sms", func(w http.ResponseWriter, r *http.Request) {
			var msg, from string
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "sms") {
				log.Printf("SMS: failed auth '%s'\n", user)
				w.Header().Set("WWW-Authenticate", `Basic realm="sms notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = r.ParseForm()
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}

			/*
				POST /_sms HTTP/1.0
				X-Forwarded-For: 54.162.174.232
				Host: suah.dev:443
				X-Forwarded-Proto: https
				X-Forwarded-Ssl: on
				Connection: close
				Content-Length: 407
				Content-Type: application/x-www-form-urlencoded
				X-Twilio-Signature: Y0YhchuDX0NZqhN7RmrdWWjg/mM=
				I-Twilio-Idempotency-Token: 10152b5f-8f76-4477-95e8-ed2a7000675b
				User-Agent: TwilioProxy/1.1

				ToCountry=US&ToState=MI&SmsMessageSid=SM0ba92af93fd5312949198953b503a787&NumMedia=0&ToCity=&FromZip=80919&SmsSid=SM0ba92af93fd5312949198953b503a787&FromState=CO&SmsStatus=received&FromCity=COLORADO+SPRINGS&Body=New+test&FromCountry=US&To=%2B18105103020&ToZip=&NumSegments=1&MessageSid=SM0ba92af93fd5312949198953b503a787&AccountSid=ACa42644b105e1329a42ca640daa61903b&From=%2B17192103020&ApiVersion=2010-04-01
			*/

			switch r.Method {
			case http.MethodGet, http.MethodPost:
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

			msg = strings.TrimSuffix(msg, "\n")

			if msg == "" {
				fmt.Fprintf(w, "empty message")
				return
			}

			for _, p := range *plugins {
				if p.Match(from, msg) {
					log.Printf("%s: responding to '%s'", p.Name(), from)
					p.SetStore(store)

					start := time.Now()
					resp, err := p.Process(from, msg)
					if err != nil {
						fmt.Println(err)
					}

					log.Println(resp)

					elapsed := time.Since(start)
					if verbose {
						log.Printf("%s took %s to run\n", p.Name(), elapsed)
					}
				}

			}
			/*
				for _, line := range strings.Split(msg, "\n") {
					log.Printf("SMS: sending '%s'\n", line)
					err = plugins.SendUnescNotice(cli, smsRoom, line)
					if err != nil {
						http.Error(
							w,
							fmt.Sprintf("can not send commit info: %s", err),
							http.StatusInternalServerError,
						)
						return
					}
				}
			*/

			fmt.Fprintf(w, "ok")

		})

		log.Fatal(http.ListenAndServe(smsPort, nil))
	}
}
