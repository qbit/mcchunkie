package plugins

import (
	"github.com/matrix-org/gomatrix"
	"math/rand"
	"regexp"
	"time"
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

func (h *Groan) resp() string {
	a := []string{
		"Ugh.",
		"ugh",
		"ffffuuuu",
		"sigh.",
		"oh fml.",
		"........",
	}

	rand.Seed(time.Now().Unix())
	return a[rand.Intn(len(a))]

}

// RespondText to groan events
func (h *Groan) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendText(c, ev.RoomID, h.resp())
}

// Name returns the name of the Groan plugin
func (h *Groan) Name() string {
	return "Groan"
}
