package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// Covid responds to covid requests
type Covid struct {
}

// State represents a individual state from the json api
type State struct {
	Confirmed int `json:"confirmed"`
	Recovered int `json:"recovered"`
	Deaths    int `json:"deaths"`
}

// Descr describes this plugin
func (h *Covid) Descr() string {
	return "Queries [thebigboard.cc](http://www.thebigboard.cc)'s api for information on COVID-19."
}

// Re returns the beer matching string
func (h *Covid) Re() string {
	return `(?i)^covid: (.+)$`
}

func (h *Covid) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// Match determines if we should call the response for Covid
func (h *Covid) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (h *Covid) SetStore(_ PluginStore) {}

// Process does the heavy lifting
func (h *Covid) Process(from, msg string) string {
	state := h.fix(msg)
	if state != "" {
		var states = make(map[string]State)
		req := HTTPRequest{
			Method:  "GET",
			ResBody: &states,
			URL:     "http://www.thebigboard.cc/feeds/v1/us.json",
		}
		_ = req.DoJSON()
		// updated and source cause some issues here
		//if err != nil {
		//	SendText(c, ev.RoomID, fmt.Sprintf("Computer says no: %s", err))
		//}

		var s State
		for i, p := range states {
			if strings.EqualFold(i, state) {
				s = p
				state = i
			}
		}
		return fmt.Sprintf("_%s_: confirmed cases: **%d**, recovered: _%d_, deaths: _%d_", state, s.Confirmed, s.Recovered, s.Deaths)
	}
	return fmt.Sprintf("invalid state: %q", state)
}

// RespondText to looking up of beer requests
func (h *Covid) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendMD(c, ev.RoomID, h.Process(ev.Sender, post))

}

// Name Covid!
func (h *Covid) Name() string {
	return "Covid"
}
