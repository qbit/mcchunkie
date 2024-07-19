package chats

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

type voipms struct {
	did         string
	dst         string
	message     string
	method      string
	apiUser     string
	apiPassword string
}

func (v voipms) UrlStr() string {
	base, _ := url.Parse("https://voip.ms/api/v1/rest.php")

	params := url.Values{}
	params.Add("did", v.did)
	params.Add("dst", v.dst)
	params.Add("method", v.method)
	params.Add("api_username", v.apiUser)
	params.Add("api_password", v.apiPassword)
	params.Add("message", v.message)

	base.RawQuery = params.Encode()
	return base.String()
}

func sendVoipmsResp(v voipms) error {
	resp, err := http.Get(v.UrlStr())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(str))
	return nil
}

// SMSListen listens for our incoming sms
func SMSListen(store ChatStore, plugins *plugins.Plugins) {
	var (
		smsPort, _    = store.Get("sms_listen")
		smsAllowed, _ = store.Get("sms_users")
		smsUsers      = strings.Split(smsAllowed, ",")
		voipmsUser, _ = store.Get("voipms_user")
		voipmsPass, _ = store.Get("voipms_api_pass")
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
				// voip.ms
				// to={TO}&from={FROM}&message={MESSAGE}&id={ID}&date={TIMESTAMP}
				msg = r.URL.Query().Get("message")
				from = r.URL.Query().Get("from")
				to := r.URL.Query().Get("to")

				log.Println("dst:", to, "did:", from)

				for _, p := range *plugins {
					if p.Match(from, msg) {
						log.Printf("%s: responding to '%s'", p.Name(), from)
						p.SetStore(store)

						resp := p.Process(from, msg)
						err := sendVoipmsResp(voipms{
							did:         to,
							dst:         from,
							message:     resp,
							method:      "sendSMS",
							apiUser:     voipmsUser,
							apiPassword: voipmsPass,
						})
						if err != nil {
							log.Println(err)
						}
					}
				}
				return
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
