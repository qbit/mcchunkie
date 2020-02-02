package plugins

import (
	"fmt"
	"log"
	"regexp"
	"runtime"

	"github.com/matrix-org/gomatrix"
)

// Version responds to hi messages
type Version struct {
}

func (v *Version) match(msg string) bool {
	re := regexp.MustCompile(`(?i)version$`)
	return re.MatchString(msg)
}

func (v *Version) print(to string) string {
	return fmt.Sprintf("%s, I am written in Go, running on %s", to, runtime.GOOS)
}

// Respond to version events
func (v *Version) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				u := NameRE.ReplaceAllString(user, "$1")
				s := NameRE.ReplaceAllString(ev.Sender, "$1")
				if ToMe(u, post) {
					if v.match(post) {
						log.Printf("%s: responding to '%s'", v.Name(), ev.Sender)
						SendMessage(c, ev.RoomID, v.print(s))
					}
				}
			}
		}
	}
}

// Name Version
func (v *Version) Name() string {
	return "Version"
}
