package plugins

import (
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// LoveYou responds to love messages
type LoveYou struct {
}

// Match checks for 'i love you' and a reference to the bot name
func (h *LoveYou) Match(user, msg string) bool {
	re := regexp.MustCompile(`(?i)i love you`)
	return re.MatchString(msg) && ToMe(user, msg)
}

func (h *LoveYou) resp() string {
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
func (h *LoveYou) SetStore(s PluginStore) {}

// RespondText to love events
func (h *LoveYou) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
	SendText(c, ev.RoomID, h.resp())
}

// Name i love you
func (h *LoveYou) Name() string {
	return "LoveYou"
}
