package plugins

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
)

// Ban responds to ban messages
type Ban struct {
}

// Descr returns a description
func (h *Ban) Descr() string {
	return "Ban a list of users or servers."
}

// Re matches "ban:" in a given string
func (h *Ban) Re() string {
	return `(?i)^ban (user|server): (.*)$`
}

// Match determines if we should execute Ban
func (h *Ban) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store, so just return
func (h *Ban) SetStore(_ PluginStore) {}

// Process does the heavy lifting
func (h *Ban) Process(from, msg string) string {
	return ""
}

// RespondText to botsnack events
func (h *Ban) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) error {
	switch ev.Sender {
	case "@qbit:tapenet.org":
		speed := 5
		re := regexp.MustCompile(h.Re())
		cmd := re.ReplaceAllString(post, "$1")
		bans := strings.Split(re.ReplaceAllString(post, "$2"), " ")

		go func() {
			SendText(c, ev.RoomID, fmt.Sprintf("Banning %d %s with %d seconds inbetween bans.", len(bans), cmd, speed))
			for _, ban := range bans {
				st := fmt.Sprintf("hammer ban ob %s %s spam", cmd, ban)
				SendText(c, ev.RoomID, st)
				time.Sleep(time.Second * time.Duration(speed))
			}
			SendText(c, ev.RoomID, "Done banning.")
		}()
	}
	return nil
}

// Name Ban
func (h *Ban) Name() string {
	return "Ban"
}
