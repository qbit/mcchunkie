package plugins

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Thanks responds to thanks
type Thanks struct {
}

// Descr describes this plugin
func (h *Thanks) Descr() string {
	return "Bots should be respectful. Respond to thanks."
}

// Re checks for various forms of thanks
func (h *Thanks) Re() string {
	return `(?i)^thank you|thank you$|^thanks|thanks$|^ty|ty$`
}

// Match determines if we are being thanked
func (h *Thanks) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Thanks) SetStore(_ PluginStore) {}

// Process
func (h *Thanks) Process(from, post string) string {
	s := NameRE.ReplaceAllString(from, "$1")
	a := []string{
		fmt.Sprintf("welcome %s", s),
		"welcome",
		"you're welcome",
		"np!",
		fmt.Sprintf("you're welcome, %s", s),
	}

	return a[rand.Intn(len(a))]
}

// RespondText to welcome back events
func (h *Thanks) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	return SendText(c, ev.RoomID, h.Process(ev.Sender, ""))
}

// Name Thanks
func (h *Thanks) Name() string {
	return "Thanks"
}
