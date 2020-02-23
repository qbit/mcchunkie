package plugins

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Snap responds to OpenBSD snapshot checks
type Snap struct {
}

// Descr describes this plugin
func (p *Snap) Descr() string {
	return "checks the current build date of OpenBSD snapshots."
}

// Re returns the federation check matching string
func (p *Snap) Re() string {
	return `(?i)^snap:$`
}

// Match determines if we should call the response for Snap
func (p *Snap) Match(user, msg string) bool {
	re := regexp.MustCompile(p.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (p *Snap) SetStore(s PluginStore) {}

// RespondText to looking up of federation check requests
func (p *Snap) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	log.Printf("%s: responding to '%s'", p.Name(), ev.Sender)
	resp, err := http.Get("https://ftp.usa.openbsd.org/pub/OpenBSD/snapshots/amd64/BUILDINFO")
	if err != nil {
		SendText(c, ev.RoomID, fmt.Sprintf("%s", err))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendText(c, ev.RoomID, fmt.Sprintf("%s", err))
		return
	}

	SendText(c, ev.RoomID, string(body))
}

// Name Snap!
func (p *Snap) Name() string {
	return "Snap"
}
