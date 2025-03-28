package plugins

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
)

// DMRUser represents a response from:
// https://database.radioid.net/api/dmr/user/
type DMRUser struct {
	Count   int `json:"count"`
	Results []struct {
		Callsign string `json:"callsign"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Fname    string `json:"fname"`
		ID       int    `json:"id"`
		Remarks  string `json:"remarks"`
		State    string `json:"state"`
		Surname  string `json:"surname"`
	} `json:"results"`
}

// DMRRepeater represents a response from:
// https://database.radioid.net/api/dmr/repeater/
type DMRRepeater struct {
	Count   int `json:"count"`
	Results []struct {
		Callsign    string `json:"callsign"`
		City        string `json:"city"`
		ColorCode   int    `json:"color_code"`
		Country     string `json:"country"`
		Details     string `json:"details"`
		Frequency   string `json:"frequency"`
		ID          string `json:"id"`
		IpscNetwork string `json:"ipsc_network"`
		Offset      string `json:"offset"`
		Rfinder     string `json:"rfinder"`
		State       string `json:"state"`
		Trustee     string `json:"trustee"`
		TsLinked    string `json:"ts_linked"`
	} `json:"results"`
}

// DMR is our plugin type
type DMR struct {
}

// SetStore is the setup function for a plugin
func (p *DMR) SetStore(s PluginStore) {
}

// Descr describes this plugin
func (p *DMR) Descr() string {
	return "Queries radioid.net"
}

// Re is what our dmr request matches
func (p *DMR) Re() string {
	return `(?i)^dmr (user|repeater) (surname|id|callsign|city|county) (.+)$`
}

// Match checks for "dmr " messages
func (p *DMR) Match(_, msg string) bool {
	re := regexp.MustCompile(p.Re())
	return re.MatchString(msg)
}

func (p *DMR) param(msg string) string {
	re := regexp.MustCompile(p.Re())
	return re.ReplaceAllString(msg, "$2")
}

func (p *DMR) mode(msg string) string {
	re := regexp.MustCompile(p.Re())
	return re.ReplaceAllString(msg, "$1")
}

func (p *DMR) query(msg string) string {
	re := regexp.MustCompile(p.Re())
	return re.ReplaceAllString(msg, "$3")
}

func (p *DMR) Process(from, post string) (string, func() string) {
	mode := p.mode(post)
	param := p.param(post)
	search := p.query(post)

	endpoint := "https://database.radioid.net/api/dmr/%s/?%s"

	params := url.Values{}
	params.Set(param, search)

	u := fmt.Sprintf(endpoint, mode, params.Encode())

	var req = HTTPRequest{
		Timeout: 10 * time.Second,
		URL:     u,
		Method:  "GET",
	}

	switch mode {
	case "repeater":
		var res = &DMRRepeater{}
		req.ResBody = res
		err := req.DoJSON()
		if err != nil {
			return err.Error(), RespStub
		}

		if res.Count == 0 {
			return fmt.Sprintf("nothing found for '%s'", params.Encode()), RespStub
		}

		var s []string
		s = append(s, fmt.Sprintf("**Callsign**: %s", res.Results[0].Callsign))
		s = append(s, fmt.Sprintf("**ID**: %s", res.Results[0].ID))
		s = append(s, fmt.Sprintf("**Frequency**: %s", res.Results[0].Frequency))
		s = append(s, fmt.Sprintf("**Offset**: %s", res.Results[0].Offset))

		return strings.Join(s, ", "), RespStub

	case "user":
		var res = &DMRUser{}
		req.ResBody = res
		err := req.DoJSON()
		if err != nil {
			return err.Error(), RespStub
		}

		if res.Count == 0 {
			return fmt.Sprintf("nothing found for '%s'", params.Encode()), RespStub
		}

		var s []string
		s = append(s, fmt.Sprintf("**Name**: %s %s", res.Results[0].Fname, res.Results[0].Surname))
		s = append(s, fmt.Sprintf("**ID**: %d", res.Results[0].ID))
		s = append(s, fmt.Sprintf("**Callsign**: %s", res.Results[0].Callsign))

		return strings.Join(s, ", "), RespStub
	}
	return fmt.Sprintf("invalid mode: %q", mode), RespStub
}

// RespondText to looking up of DMR info
func (p *DMR) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := p.Process(ev.Sender, post)
	return SendMD(c, ev.RoomID, resp)
}

// Name DMR!
func (p *DMR) Name() string {
	return "DMR"
}
