package plugins

import (
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// Yeah puts on the shades
type Yeah struct {
}

// Descr describes this plugin
func (h *Yeah) Descr() string {
	return "Now you're cool."
}

// Re are the regexes that high five uses
func (h *Yeah) Re() string {
	return `(?i)CSI$`
}

// SetStore we don't need a store here.
func (h *Yeah) SetStore(_ PluginStore) {}

// Match determines if we should bother giving a high five
func (h *Yeah) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return ToMe(user, msg) && re.MatchString(msg)
}

func (h *Yeah) Process(from, post string) string {
	return ""
}

// RespondText to high five events
func (h *Yeah) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	parts := []string{
		"( •_•)",
		"( •_•)>⌐■-■",
		"(⌐■_■)",
	}
	for _, p := range parts {
		SendText(c, ev.RoomID, p)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(5 * time.Second)
	return SendText(c, ev.RoomID, "YEEEAAAAAAHHHHHH!")
}

// Name returns the name of the Yeah plugin
func (h *Yeah) Name() string {
	return "Yeah"
}
