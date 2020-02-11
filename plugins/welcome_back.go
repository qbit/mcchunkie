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

// Descr describes this plugin
func (h *Wb) Descr() string {
	return "Respond to welcome back messages."
}

// Re checks for various welcome back things
func (h *Wb) Re() string {
	return `(?i)^welcome back|welcome back$|^wb|wb$`
}

// Match determins if we are welcomed back
func (h *Wb) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Wb) SetStore(s PluginStore) {}

// RespondText to welcome back events
func (h *Wb) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
	SendText(c, ev.RoomID, fmt.Sprintf("thanks %s!", s))
}

// Name Wb
func (h *Wb) Name() string {
	return "Wb"
}
