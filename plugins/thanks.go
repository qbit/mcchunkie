package plugins

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

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

// Match determins if we are being thanked
func (h *Thanks) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Thanks) SetStore(s PluginStore) {}

// RespondText to welcome back events
func (h *Thanks) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) error {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")
	a := []string{
		fmt.Sprintf("welcome %s", s),
		"welcome",
		"you're welcome",
		"np!",
		fmt.Sprintf("you're welcome, %s", s),
	}

	rand.Seed(time.Now().Unix())

	return SendText(c, ev.RoomID, a[rand.Intn(len(a))])
}

// Name Thanks
func (h *Thanks) Name() string {
	return "Thanks"
}
