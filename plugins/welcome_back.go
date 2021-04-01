package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Wb responds to welcome back messages
type Wb struct {
}

// Descr describes this plugin
func (h *Wb) Descr() string {
	return "Respond to welcome back messages."
}

// Re checks for various welcome back things
func (h *Wb) Re() string {
	return `(?i)^welcome back|welcome back$|^wb|wb$`
}

// Match determines if we are welcomed back
func (h *Wb) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore we don't need a store here
func (h *Wb) SetStore(_ PluginStore) {}

func (h *Wb) Process(from, post string) string {
	s := NameRE.ReplaceAllString(from, "$1")
	return fmt.Sprintf("thanks %s!", s)
}

// RespondText to welcome back events
func (h *Wb) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	return SendText(c, ev.RoomID, h.Process(ev.Sender, ""))

}

// Name Wb
func (h *Wb) Name() string {
	return "Wb"
}
