package plugins

import (
	"fmt"
	"regexp"
	"time"

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
func (h *HighFive) SetStore(s PluginStore) {}

// Match determines if we should bother giving a high five
func (h *HighFive) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return ToMe(user, msg) && re.MatchString(msg)
}

// RespondText to high five events
func (h *HighFive) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) error {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	rm := regexp.MustCompile(rightFive())
	lm := regexp.MustCompile(leftFive())

	if rm.MatchString(post) {
		_ = SendText(c, ev.RoomID, fmt.Sprintf("\\o %s", s))
		time.Sleep(time.Second * 5)
		return SendText(c, ev.RoomID, fmt.Sprintf("now go wash your hands, %s", s))
	}
	if lm.MatchString(post) {
		_ = SendText(c, ev.RoomID, fmt.Sprintf("%s o/", s))
		time.Sleep(time.Second * 5)
		return SendText(c, ev.RoomID, fmt.Sprintf("now go wash your hands, %s", s))
	}
	return nil
}

// Name returns the name of the HighFive plugin
func (h *HighFive) Name() string {
	return "HighFive"
}
