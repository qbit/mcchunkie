package plugins

import (
	"fmt"
	"log"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Source responds to source requests
type Source struct {
}

func (h *Source) match(msg string) bool {
	re := regexp.MustCompile(`(?i)where is your (source|code)`)
	return re.MatchString(msg)
}

// Respond to questions about TheSource™©®⑨
func (h *Source) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				u := NameRE.ReplaceAllString(user, "$1")
				s := NameRE.ReplaceAllString(ev.Sender, "$1")
				if ToMe(u, post) {
					if h.match(post) {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						SendText(c, ev.RoomID, fmt.Sprintf("%s: %s ;D", s, "https://git.sr.ht/~qbit/mcchunkie"))
					}
				}
			}
		}
	}
}

// Name Source
func (h *Source) Name() string {
	return "Source"
}
