package plugins

import (
	"fmt"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// Beat responds to beat messages
type Beat struct {
}

// Descr describes this plugin
func (h *Beat) Descr() string {
	return "Print the current [beat time](https://en.wikipedia.org/wiki/Swatch_Internet_Time)."
}

// Re is the regex for matching beat messages.
func (h *Beat) Re() string {
	return `(?i)^\.beat$|^what time is it[\?!]+$|^beat( )?time:?\??$`
}

// Match determines if we are asking for a beat
func (h *Beat) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *Beat) SetStore(_ PluginStore) {}

// RespondText to beat request events
func (h *Beat) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	resp, _ := h.Process("", "")
	return SendText(c, ev.RoomID, resp)
}

// Process does the heavy lifting of calculating .beat
func (h *Beat) Process(from, msg string) (string, func() string) {
	n := time.Now()
	utc1 := n.Unix() + 3600
	r := utc1 % 86400
	bt := float32(r) / 86.4
	return fmt.Sprintf("@%03d", int32(bt)), RespStub
}

// Name beat
func (h *Beat) Name() string {
	return "Beat"
}
