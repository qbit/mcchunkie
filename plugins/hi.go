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

// Match determines if we are highfiving
func (h *Hi) Match(user, msg string) bool {
	re := regexp.MustCompile(`(?i)^hi|hi$`)
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Hi) SetStore(s PluginStore) {}

// RespondText to hi events
func (h *Hi) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
	SendText(c, ev.RoomID, fmt.Sprintf("hi %s!", s))
}

// Name hi
func (h *Hi) Name() string {
	return "Hi"
}
