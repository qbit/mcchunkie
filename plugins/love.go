package plugins

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// LoveYou responds to love messages
type LoveYou struct {
}

// Descr describes this plugin
func (h *LoveYou) Descr() string {
	return "Spreading love where ever we can by responding when someone shows us love."
}

// Re matches "i love you"
func (h *LoveYou) Re() string {
	return `(?i)i love you`
}

// Match checks for 'i love you' and a reference to the bot name
func (h *LoveYou) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// Process does the heavy lifting
func (h *LoveYou) Process(from, post string) string {
	a := []string{
		"I am not ready for this kind of relationship!",
		"ಠ_ಠ",
		"I love you too!",
		"(╯‵Д′)╯彡┻━┻",
		"hawkard!",
	}

	rand.Seed(time.Now().Unix())
	return a[rand.Intn(len(a))]

}

// SetStore we don't need a store, so just return
func (h *LoveYou) SetStore(_ PluginStore) {}

// RespondText to love events
func (h *LoveYou) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	return SendText(c, ev.RoomID, h.Process("", ""))
}

// Name i love you
func (h *LoveYou) Name() string {
	return "LoveYou"
}
