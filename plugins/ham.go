package plugins

import (
	"fmt"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// LicenseResp represents a response from https://data.fcc.gov/api/license-view/basicSearch/getLicenses
type LicenseResp struct {
	Status   string   `json:"status"`
	Licenses Licenses `json:"Licenses"`
}

// License is an individual license
type License struct {
	LicName      string `json:"licName"`
	Frn          string `json:"frn"`
	Callsign     string `json:"callsign"`
	CategoryDesc string `json:"categoryDesc"`
	ServiceDesc  string `json:"serviceDesc"`
	StatusDesc   string `json:"statusDesc"`
	ExpiredDate  string `json:"expiredDate"`
	LicenseID    string `json:"licenseID"`
	LicDetailURL string `json:"licDetailURL"`
}

// Licenses is a collection of individual licenses.
type Licenses struct {
	Page       string    `json:"page"`
	RowPerPage string    `json:"rowPerPage"`
	TotalRows  string    `json:"totalRows"`
	LastUpdate string    `json:"lastUpdate"`
	License    []License `json:"License"`
}

// Ham for querying the fcc'd uls
type Ham struct{}

// Descr describes this plugin
func (h *Ham) Descr() string {
	return "queries the FCC's [ULS](https://wireless2.fcc.gov/UlsApp/UlsSearch/searchLicense.jsp) for a given callsign."
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
func (h *Ham) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (h *Ham) SetStore(s PluginStore) {}

func (h *Ham) pretty(resp *LicenseResp) string {
	return fmt.Sprintf("%s: %s (expires: %s) %s\n",
		resp.Licenses.License[0].Callsign,
		resp.Licenses.License[0].LicName,
		resp.Licenses.License[0].ExpiredDate,
		resp.Licenses.License[0].CategoryDesc,
	)
}

// RespondText to looking up of federation check requests
func (h *Ham) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	call := h.fix(post)
	if call != "" {
		furl := fmt.Sprintf("%s%s",
			"http://data.fcc.gov/api/license-view/basicSearch/getLicenses?format=json&searchValue=",
			call,
		)

		var res = &LicenseResp{}
		var req = HTTPRequest{
			Timeout: 10 * time.Second,
			URL:     furl,
			Method:  "GET",
			ResBody: res,
		}

		err := req.DoJSON()
		if err != nil {
			SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look things up in ULS (%s)", ev.Sender, err))
			return
		}

		if res.Status == "OK" {
			SendText(c, ev.RoomID, h.pretty(res))
		} else {
			SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look things up in ULS. The response was not OK.", ev.Sender))
		}
	}
}

// Name Ham!
func (h *Ham) Name() string {
	return "Ham"
}
