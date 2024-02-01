package plugins

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

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
func (p *Snap) Match(_, msg string) bool {
	re := regexp.MustCompile(p.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (p *Snap) SetStore(_ PluginStore) {}

// Process does the heavy lifting
func (p *Snap) Process(from, post string) string {
	snapResp, err := http.Get("https://ftp.usa.openbsd.org/pub/OpenBSD/snapshots/amd64/BUILDINFO")
	if err != nil {
		return err.Error()
	}
	defer snapResp.Body.Close()

	buildBody, err := io.ReadAll(snapResp.Body)
	if err != nil {
		return err.Error()
	}

	str := string(buildBody)
	parts := strings.Split(str, " - ")
	if len(parts) != 2 {
		return "Can't parse BUILDINFO"
	}

	snapDate, err := time.Parse(time.UnixDate, strings.TrimSpace(parts[1]))
	if err != nil {
		return err.Error()
	}

	pkgResp, err := http.Get("https://ftp3.usa.openbsd.org/pub/OpenBSD/snapshots/packages/amd64/SHA256")
	if err != nil {
		return err.Error()
	}
	defer pkgResp.Body.Close()

	lm := strings.TrimSpace(pkgResp.Header.Get("last-modified"))
	if lm == "" {
		return "Missing last-modified for SHA256"
	}

	pkgDate, err := time.Parse(time.RFC1123, lm)
	if lm == "" {
		return err.Error()
	}

	if pkgDate.Before(snapDate) {
		return "🔴: packages are behind snapshots. It is likely not safe to update currently!"
	}

	return "🟢:' It is safe to update your snapshot and packages!"
}

// RespondText to looking up of federation check requests
func (p *Snap) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	return SendText(c, ev.RoomID, p.Process("", ""))
}

// Name Snap!
func (p *Snap) Name() string {
	return "Snap"
}
