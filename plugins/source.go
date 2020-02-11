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

// Match determins if someone is asking about the source code
func (h *Source) Match(user, msg string) bool {
	re := regexp.MustCompile(`(?i)where is your (source|code)`)
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore does nothing in here
func (h *Source) SetStore(s PluginStore) {}

// RespondText to questions about TheSource™©®⑨
func (h *Source) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
	SendText(c, ev.RoomID, fmt.Sprintf("%s: %s ;D", s, "https://git.sr.ht/~qbit/mcchunkie"))
}

// Name Source
func (h *Source) Name() string {
	return "Source"
}
