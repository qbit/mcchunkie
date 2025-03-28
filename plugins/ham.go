package plugins

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
)

// LicenseResp represents a response from http://hamdb.org/api
type LicenseResp struct {
	Hamdb struct {
		Version  string `json:"version"`
		Callsign struct {
			Call    string `json:"call"`
			Class   string `json:"class"`
			Expires string `json:"expires"`
			Status  string `json:"status"`
			Grid    string `json:"grid"`
			Lat     string `json:"lat"`
			Lon     string `json:"lon"`
			Fname   string `json:"fname"`
			Mi      string `json:"mi"`
			Name    string `json:"name"`
			Suffix  string `json:"suffix"`
			Addr1   string `json:"addr1"`
			Addr2   string `json:"addr2"`
			State   string `json:"state"`
			Zip     string `json:"zip"`
			Country string `json:"country"`
		} `json:"callsign"`
		Messages struct {
			Status string `json:"status"`
		} `json:"messages"`
	} `json:"hamdb"`
}

// Ham for querying the fcc'd uls
type Ham struct{}

// Descr describes this plugin
func (h *Ham) Descr() string {
	return "queries HamDB.org for a given callsign."
}

// Re returns the federation check matching string
func (h *Ham) Re() string {
	return `(?i)^ham: (\w+)$`
}

func (h *Ham) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// Match determines if we should call the response for Ham
func (h *Ham) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (h *Ham) SetStore(_ PluginStore) {}

func (h *Ham) pretty(resp *LicenseResp) string {
	var s []string

	s = append(s, fmt.Sprintf("%s: %s %s (expires: %s) %s (%s)\n",
		resp.Hamdb.Callsign.Call,
		resp.Hamdb.Callsign.Fname, resp.Hamdb.Callsign.Name,
		resp.Hamdb.Callsign.Expires,
		resp.Hamdb.Callsign.Country,
		resp.Hamdb.Callsign.Grid,
	))

	return strings.Join(s, " ")
}

// Process does the heavy lifting
func (h *Ham) Process(from, post string) (string, func() string) {
	call := h.fix(post)
	if call != "" {
		furl := fmt.Sprintf("http://api.hamdb.org/v1/%s/json/mcchunkie",
			call,
		)

		var res = &LicenseResp{}
		var req = HTTPRequest{
			Timeout: 15 * time.Second,
			URL:     furl,
			Method:  "GET",
			ResBody: res,
		}

		err := req.DoJSON()
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look things up in ULS (%s)", from, err), RespStub
		}

		if res.Hamdb.Messages.Status == "OK" {
			return h.pretty(res), RespStub
		}

		return fmt.Sprintf("sorry %s, I can't look things up in ULS. The response was not OK.", from), RespStub
	}

	return "invalid callsign", RespStub
}

// RespondText to looking up of federation check requests
func (h *Ham) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name Ham!
func (h *Ham) Name() string {
	return "Ham"
}
