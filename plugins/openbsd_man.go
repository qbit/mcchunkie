package plugins

import (
	"fmt"
	"log"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// OpenBSDMan responds to beer requests
type OpenBSDMan struct {
}

func (h *OpenBSDMan) re() string {
	return `(?i)^man: ([1-9]?p?)\s?(\w+)$`
}

func (h *OpenBSDMan) fix(msg string) string {
	re := regexp.MustCompile(h.re())
	resp := ""
	section := re.ReplaceAllString(msg, "$1")
	if section == msg {
		return ""
	}
	if section != "" {
		resp = re.ReplaceAllString(msg, "$2.$1")
	} else {
		resp = re.ReplaceAllString(msg, "$2")
	}

	return resp
}

func (h *OpenBSDMan) match(msg string) bool {
	re := regexp.MustCompile(h.re())
	return re.MatchString(msg)
}

// RespondText sends back a man page.
func (h *OpenBSDMan) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	if h.match(post) {
		page := h.fix(post)
		if page != "" {
			log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
			SendText(c, ev.RoomID, fmt.Sprintf("https://man.openbsd.org/%s", page))
		}
	}
}

// Name OpenBSDMan!
func (h *OpenBSDMan) Name() string {
	return "OpenBSDMan"
}
