package plugins

import (
	"log"
	"math/rand"
	"regexp"
	"time"

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
func (h *BotSnack) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *BotSnack) resp() string {
	a := []string{
		"omm nom nom nom",
		"*puke*",
		"MOAR!",
		"=.=",
	}

	rand.Seed(time.Now().Unix())
	return a[rand.Intn(len(a))]

}

// SetStore we don't need a store, so just return
func (h *BotSnack) SetStore(s PluginStore) {}

// RespondText to botsnack events
func (h *BotSnack) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	u := NameRE.ReplaceAllString(user, "$1")
	if ToMe(u, post) {
		log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
		SendText(c, ev.RoomID, h.resp())
	}
}

// Name BotSnack
func (h *BotSnack) Name() string {
	return "BotSnack"
}
