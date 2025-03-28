package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Source responds to source requests
type Source struct {
}

// Descr describes this plugin
func (h *Source) Descr() string {
	return "Tell people where they can find more information about myself."
}

// Re matches the source code question
func (h *Source) Re() string {
	return `(?i)where is your (source|code)`
}

// Match determines if someone is asking about the source code
func (h *Source) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

// SetStore does nothing in here
func (h *Source) SetStore(_ PluginStore) {}

// Process does the heavy lifting
func (h *Source) Process(from, post string) (string, func() string) {
	s := NameRE.ReplaceAllString(from, "$1")
	return fmt.Sprintf("%s: %s ;D", s, "https://git.sr.ht/~qbit/mcchunkie"), RespStub
}

// RespondText to questions about TheSource™©®⑨
func (h *Source) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	resp, _ := h.Process(ev.Sender, "")
	return SendText(c, ev.RoomID, resp)
}

// Name Source
func (h *Source) Name() string {
	return "Source"
}
