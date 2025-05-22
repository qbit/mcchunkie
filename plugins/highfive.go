package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// HighFive high fives!
type HighFive struct {
}

func rightFive() string {
	return "o/"
}

func leftFive() string {
	return `\\o`
}

// Descr describes this plugin
func (h *HighFive) Descr() string {
	return "Everyone loves highfives."
}

// Re are the regexes that high five uses
func (h *HighFive) Re() string {
	return fmt.Sprintf("%s|%s", rightFive(), leftFive())
}

// SetStore we don't need a store here.
func (h *HighFive) SetStore(_ PluginStore) {}

// Match determines if we should bother giving a high five
func (h *HighFive) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return ToMe(user, msg) && re.MatchString(msg)
}

func (h *HighFive) Process(from, post string) (string, func() string) {
	s := NameRE.ReplaceAllString(from, "$1")

	rm := regexp.MustCompile(rightFive())
	lm := regexp.MustCompile(leftFive())

	if rm.MatchString(post) {
		return fmt.Sprintf("\\o %s", s), func() string { return "" }
	}

	if lm.MatchString(post) {
		return fmt.Sprintf("%s o/", s), func() string { return "" }
	}

	return "\\o/", func() string { return "" }
}

// RespondText to high five events
func (h *HighFive) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name returns the name of the HighFive plugin
func (h *HighFive) Name() string {
	return "HighFive"
}
