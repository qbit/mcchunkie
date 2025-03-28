package plugins

import (
	"fmt"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

type BananaStab struct {
}

// Descr describes this plugin
func (h *BananaStab) Descr() string {
	return "Stab someone with some bananas."
}

// Re matches our stabbing format
func (h *BananaStab) Re() string {
	return `(?i)^stab: (.+)$`
}

func (h *BananaStab) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	stabee := re.ReplaceAllString(msg, "$1")
	return stabee
}

// Match checks for our stabee person
func (h *BananaStab) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore does nothing in BananaStab
func (h *BananaStab) SetStore(_ PluginStore) {}

func (h *BananaStab) Process(from, post string) (string, func() string) {
	stabee := h.fix(post)
	stabtxt := "..."
	if stabee != "" {
		stabtxt = fmt.Sprintf("stabs %s with the fury of a thousand radioactive bananas", stabee)
	}
	//jsonmsg := "{ \"body\": \"" + stabtxt + "\", \"type\": \"m.emote\"}"
	return stabtxt, RespStub
}

// RespondText stabs an unsuspecting person
func (h *BananaStab) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendEmote(c, ev.RoomID, resp)
}

// Name BananaStab!
func (h *BananaStab) Name() string {
	return "BananaStab"
}
