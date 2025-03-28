package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// OpenBSDMan responds to beer requests
type OpenBSDMan struct {
}

// Descr describes this plugin
func (h *OpenBSDMan) Descr() string {
	return "Produces a link to man.openbsd.org."
}

// Re matches our man format
func (h *OpenBSDMan) Re() string {
	return `(?i)^man: ([1-9][p]?)?\s?(.+)$`
}

func (h *OpenBSDMan) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	resp := ""
	section := re.ReplaceAllString(msg, "$1")
	if section == msg {
		return ""
	}
	if section != "" {
		resp = re.ReplaceAllString(msg, "$2.$1")
		if matched, _ := regexp.MatchString(`3p`, resp); matched {
			resp = fmt.Sprintf("man3p/%s", resp)
		}

	} else {
		resp = re.ReplaceAllString(msg, "$2")
	}

	return resp
}

// Match checks for our man page re
func (h *OpenBSDMan) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore does nothing in OpenBSDMan
func (h *OpenBSDMan) SetStore(_ PluginStore) {}

func (h *OpenBSDMan) Process(from, post string) (string, func() string) {
	page := h.fix(post)
	if page != "" {
		return fmt.Sprintf("https://man.openbsd.org/%s", page), RespStub
	}
	return "...", RespStub
}

// RespondText sends back a man page.
func (h *OpenBSDMan) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, delayedResp := h.Process(ev.Sender, post)
	go func() {
		SendText(c, ev.RoomID, delayedResp())
	}()

	return SendText(c, ev.RoomID, resp)
}

// Name OpenBSDMan!
func (h *OpenBSDMan) Name() string {
	return "OpenBSDMan"
}
