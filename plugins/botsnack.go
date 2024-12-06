package plugins

import (
	"math/rand"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// BotSnack responds to botsnack messages
type BotSnack struct {
}

// Descr returns a description
func (h *BotSnack) Descr() string {
	return "Consumes a botsnack. This pleases mcchunkie and brings balance to the universe."
}

// Re matches "botsnack" in a given string
func (h *BotSnack) Re() string {
	return `(?i)botsnack`
}

// Match determines if we should execute BotSnack
func (h *BotSnack) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store, so just return
func (h *BotSnack) SetStore(_ PluginStore) {}

// Process does the heavy lifting
func (h *BotSnack) Process(from, msg string) string {
	u := NameRE.ReplaceAllString(from, "$1")
	if ToMe(u, msg) {
		a := []string{
			"omm nom nom nom",
			"*puke*",
			"MOAR!",
			"=.=",
		}

		return a[rand.Intn(len(a))]
	}
	return ""
}

// RespondText to botsnack events
func (h *BotSnack) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) error {
	return SendText(c, ev.RoomID, h.Process(ev.Sender, post))
}

// Name BotSnack
func (h *BotSnack) Name() string {
	return "BotSnack"
}
