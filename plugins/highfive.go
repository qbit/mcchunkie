package plugins

import (
	"fmt"
	"log"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// HighFive high fives!
type HighFive struct {
}

// Respond to high five events
func (h *HighFive) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				u := NameRE.ReplaceAllString(user, "$1")
				s := NameRE.ReplaceAllString(ev.Sender, "$1")
				if ToMe(u, post) {
					if strings.Contains(post, "o/") {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						SendMessage(c, ev.RoomID, fmt.Sprintf("\\o %s", s))
					}
					if strings.Contains(post, "\\o") {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						SendMessage(c, ev.RoomID, fmt.Sprintf("%s o/", s))
					}
				}
			}
		}
	}
}

// Name returns the name of the HighFive plugin
func (h *HighFive) Name() string {
	return "HighFive"
}
