package plugins

import (
	"fmt"
	"log"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Wb responds to welcome back messages
type Wb struct {
}

func (h *Wb) match(msg string) bool {
	re := regexp.MustCompile(`(?i)^welcome back|welcome back$|^wb|wb$`)
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *Wb) SetStore(s PluginStore) { return }

// RespondText to welcome back events
func (h *Wb) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	u := NameRE.ReplaceAllString(user, "$1")
	s := NameRE.ReplaceAllString(ev.Sender, "$1")
	if ToMe(u, post) {
		if h.match(post) {
			log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
			SendText(c, ev.RoomID, fmt.Sprintf("thanks %s!", s))
		}
	}
}

// Name Wb
func (h *Wb) Name() string {
	return "Wb"
}
