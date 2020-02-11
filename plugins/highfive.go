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

// SetStore we don't need a store here.
func (h *HighFive) SetStore(s PluginStore) {}

// Match determines if we should bother giving a high five
func (h *HighFive) Match(user, msg string) bool {
	return ToMe(user, msg)
}

// RespondText to high five events
func (h *HighFive) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	if strings.Contains(post, "o/") {
		log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
		SendText(c, ev.RoomID, fmt.Sprintf("\\o %s", s))
	}
	if strings.Contains(post, "\\o") {
		log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
		SendText(c, ev.RoomID, fmt.Sprintf("%s o/", s))
	}
}

// Name returns the name of the HighFive plugin
func (h *HighFive) Name() string {
	return "HighFive"
}
