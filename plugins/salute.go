package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Salute high fives!
type Salute struct {
}

func rightSalute() string {
	return "o7"
}

// Descr describes this plugin
func (h *Salute) Descr() string {
	return "Everyone loves Salutes."
}

// Re are the regexes that high five uses
func (h *Salute) Re() string {
	return rightSalute()
}

// SetStore we don't need a store here.
func (h *Salute) SetStore(_ PluginStore) {}

// Match determines if we should bother giving a salute
func (h *Salute) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return ToMe(user, msg) && re.MatchString(msg)
}

func (h *Salute) Process(from, post string) (string, func() string) {
	s := NameRE.ReplaceAllString(from, "$1")

	rm := regexp.MustCompile(rightSalute())

	if rm.MatchString(post) {
		return fmt.Sprintf("%s o7", s), RespStub
	}

	return "o7", RespStub
}

// RespondText to high five events
func (h *Salute) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name returns the name of the Salute plugin
func (h *Salute) Name() string {
	return "Salute"
}
