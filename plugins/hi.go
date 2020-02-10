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
	re := regexp.MustCompile(`(?i)^hi|hi$`)
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *Hi) SetStore(s PluginStore) {}

// RespondText to hi events
func (h *Hi) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	u := NameRE.ReplaceAllString(user, "$1")
	s := NameRE.ReplaceAllString(ev.Sender, "$1")
	if ToMe(u, post) {
		if h.match(post) {
			log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
			SendText(c, ev.RoomID, fmt.Sprintf("hi %s!", s))
		}
	}
}

// Name hi
func (h *Hi) Name() string {
	return "Hi"
}
