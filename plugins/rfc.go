package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// RFC sends rfc urls when someone references an rfc 
type RFC struct {
}

// Descr describes this plugin
func (h *RFC) Descr() string {
	return "Produces a link to tools.ietf.org."
}

// Re matches our man format
func (h *RFC) Re() string {
	return `(?i)^rfc\s?([0-9]+)$`
}

// Match checks for our man page re
func (h *RFC) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore does nothing in RFC
func (h *RFC) SetStore(_ PluginStore) {}

// RespondText sends back a man page.
func (h *RFC) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	re := regexp.MustCompile(h.Re())
	rfcNum := re.ReplaceAllString(post, "$1")
	if rfcNum != "" {
		return SendText(c, ev.RoomID, fmt.Sprintf("https://tools.ietf.org/html/rfc%s", rfcNum))
	}

	return nil
}

// Name RFC
func (h *RFC) Name() string {
	return "RFC"
}
