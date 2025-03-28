package plugins

import (
	"math/rand"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Groan responds to groans.
type Groan struct {
}

// Descr describes this plugin
func (h *Groan) Descr() string {
	return "Ugh."
}

// Re are the regexes that high five uses
func (h *Groan) Re() string {
	return `(?i)^@groan$`
}

// SetStore we don't need a store here.
func (h *Groan) SetStore(_ PluginStore) {}

// Match determines if we should bother groaning
func (h *Groan) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Groan) Process(_, _ string) (string, func() string) {
	a := []string{
		"Ugh.",
		"ugh",
		"ffffuuuu",
		"sigh.",
		"oh fml.",
		"........",
	}

	return a[rand.Intn(len(a))], RespStub
}

// RespondText to groan events
func (h *Groan) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process("", "")
	return SendText(c, ev.RoomID, resp)
}

// Name returns the name of the Groan plugin
func (h *Groan) Name() string {
	return "Groan"
}
