package plugins

import (
	"fmt"
	"log"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Hi responds to hi messages
type Hi struct {
}

func (h *Hi) match(msg string) bool {
	re := regexp.MustCompile(`^hi|hi$`)
	return re.MatchString(msg)
}

// Respond to hi events
func (h *Hi) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				u := NameRE.ReplaceAllString(user, "$1")
				s := NameRE.ReplaceAllString(ev.Sender, "$1")
				if ToMe(u, post) {
					if h.match(post) {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						SendMessage(c, ev.RoomID, fmt.Sprintf("hi %s!", s))
					}
				}
			}
		}
	}
}

// Name hi
func (h *Hi) Name() string {
	return "Hi"
}
