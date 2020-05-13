package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Hi responds to hi messages
type Hi struct {
}

// Descr describes this plugin
func (h *Hi) Descr() string {
	return "Friendly bots say hi."
}

// Re is the regex for matching hi messages.
func (h *Hi) Re() string {
	return `(?i)^hi|hi$`
}

// Match determines if we are highfiving
func (h *Hi) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Hi) SetStore(_ PluginStore) {}

// RespondText to hi events
func (h *Hi) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	return SendText(c, ev.RoomID, fmt.Sprintf("hi %s!", s))
}

// Name hi
func (h *Hi) Name() string {
	return "Hi"
}
