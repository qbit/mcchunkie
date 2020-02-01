package plugins

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Beer responds to hi messages
type Beer struct {
}

func (h *Beer) match(msg string) string {
	re := regexp.MustCompile(`(?i)^beer: `)
	return re.ReplaceAllString(msg, "$1")
}

func (h *Beer) get(beer string) (string, error) {
	u := "https://data.opendatasoft.com/api/records/1.0/search?dataset=open-beer-database%40public-us&q="
	u = fmt.Sprintf("%s%s", u, url.PathEscape(beer))
	log.Println(u)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Respond to hi events
func (h *Beer) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				beer := h.match(post)
				if beer != "" {
					log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
					json, err := h.get(beer)
					if err != nil {
						SendMessage(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look for beer. (%s)", ev.Sender, err))
					}
					SendMessage(c, ev.RoomID, fmt.Sprintf("%s!", json))
				}
			}
		}
	}
}

// Name hi
func (h *Beer) Name() string {
	return "Beer"
}
