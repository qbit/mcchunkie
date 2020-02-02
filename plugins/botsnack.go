package plugins

import (
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// BotSnack responds to hi messages
type BotSnack struct {
}

func (h *BotSnack) match(msg string) bool {
	re := regexp.MustCompile(`(?i)botsnack`)
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

// Respond to hi events
func (h *BotSnack) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				u := NameRE.ReplaceAllString(user, "$1")
				if ToMe(u, post) {
					if h.match(post) {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						SendText(c, ev.RoomID, h.resp())
					}
				}
			}
		}
	}
}

// Name hi
func (h *BotSnack) Name() string {
	return "BotSnack"
}
