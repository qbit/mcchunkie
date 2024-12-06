package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// Help responds to hi messages
type Help struct {
}

// Descr describes this plugin
func (h *Help) Descr() string {
	return "Prints help info"
}

// Re is the regex for matching hi messages.
func (h *Help) Re() string {
	return `(?i)^help: (\w+)$`
}

// Match determines if we are highfiving
func (h *Help) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Help) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// SetStore we don't need a store here
func (h *Help) SetStore(_ PluginStore) {}

// Process does the lifting
func (h *Help) Process(from, post string) string {
	item := h.fix(post)

	var pnames []string
	for _, plg := range Plugs {
		if strings.ToLower(plg.Name()) == strings.ToLower(item) {
			return fmt.Sprintf("**%s**: `%s` -  _%s_\n", plg.Name(), plg.Re(), plg.Descr())
		}
		pnames = append(pnames, plg.Name())
	}
	return fmt.Sprintf("no help found for %q, available: %s", item, strings.Join(pnames, ", "))
}

// RespondText to hi events
func (h *Help) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendText(c, ev.RoomID, h.Process(ev.Sender, post))
}

// Name hi
func (h *Help) Name() string {
	return "Help"
}
