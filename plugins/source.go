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

// SetStore does nothing in here
func (h *Source) SetStore(s PluginStore) { return }

// RespondText to questions about TheSource™©®⑨
func (h *Source) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	u := NameRE.ReplaceAllString(user, "$1")
	s := NameRE.ReplaceAllString(ev.Sender, "$1")
	if ToMe(u, post) {
		if h.match(post) {
			log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
			SendText(c, ev.RoomID, fmt.Sprintf("%s: %s ;D", s, "https://git.sr.ht/~qbit/mcchunkie"))
		}
	}
}

// Name Source
func (h *Source) Name() string {
	return "Source"
}
