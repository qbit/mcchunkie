package plugins

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// ServiceInfo represents the version info from a response
type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Error   string `json:"error"`
}

// FedResp represents a federation statuse response
type FedResp struct {
	Status bool        `json:"FederationOK"`
	Info   ServiceInfo `json:"Version"`
}

// Feder responds to federation check requests
type Feder struct {
}

// Descr describes this plugin
func (h *Feder) Descr() string {
	return "check the Matrix federation status of a given URL."
}

// Re returns the federation check matching string
func (h *Feder) Re() string {
	return `(?i)^feder: (.*)$`
}

func (h *Feder) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// Match determines if we should call the response for Feder
func (h *Feder) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (h *Feder) SetStore(s PluginStore) {}

// RespondText to looking up of federation check requests
func (h *Feder) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	homeServer := h.fix(post)
	if homeServer != "" {
		u, err := url.Parse(fmt.Sprintf("https://%s", homeServer))
		if err != nil {
			SendText(c, ev.RoomID, "that's not a real host name.")
			return
		}

		homeServer = u.Hostname()

		furl := fmt.Sprintf("%s%s",
			"https://federationtester.matrix.org/api/report?server_name=",
			url.PathEscape(homeServer),
		)
		var fed = &FedResp{}

		var req = HTTPRequest{
			Timeout: 5 * time.Second,
			URL:     furl,
			Method:  "GET",
			ResBody: fed,
		}
		err = req.DoJSON()

		if err != nil {
			SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look up the federation status (%s)", ev.Sender, err))
			return
		}

		stat := "broken"
		if fed.Status {
			stat = "OK"
		}

		if fed.Info.Error != "" {
			SendText(c, ev.RoomID, fmt.Sprintf("%s seems to be broken, maybe it isn't a homeserver?\n%s",
				homeServer, fed.Info.Error))
		} else {
			SendText(c, ev.RoomID, fmt.Sprintf("%s is running %s (%s) and is %s.",
				homeServer, fed.Info.Name, fed.Info.Version, stat))
		}

	}
}

// Name Feder!
func (h *Feder) Name() string {
	return "Feder"
}