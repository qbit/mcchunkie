package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (h *Feder) get(hserver string) (*FedResp, error) {
	u := "https://federationtester.matrix.org/api/report?server_name="
	u = fmt.Sprintf("%s%s", u, url.PathEscape(hserver))
	hClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := hClient.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var fresp = &FedResp{}
	err = json.Unmarshal([]byte(body), fresp)
	if err != nil {
		return nil, err
	}

	return fresp, nil
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

		log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
		fed, err := h.get(homeServer)
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
