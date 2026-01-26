package plugins

import (
	"encoding/base32"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// SimpleResp is a JSON response from OpenSimpleMap.org
type SimpleResp struct {
	CreationDate      string `json:"creation_date"`
	CreationTimestamp int    `json:"creation_timestamp"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Enabled           bool   `json:"enabled"`
	ID                int    `json:"id"`
	Mailbox           struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	} `json:"mailbox"`
	Mailboxes []struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	} `json:"mailboxes"`
	LatestActivity struct {
		Action  string `json:"action"`
		Contact struct {
			Email        string `json:"email"`
			Name         any    `json:"name"`
			ReverseAlias string `json:"reverse_alias"`
		} `json:"contact"`
		Timestamp int `json:"timestamp"`
	} `json:"latest_activity"`
	NbBlock   int  `json:"nb_block"`
	NbForward int  `json:"nb_forward"`
	NbReply   int  `json:"nb_reply"`
	Note      any  `json:"note"`
	Pinned    bool `json:"pinned"`
}

// Simple is our plugin type
type Simple struct {
	db PluginStore
}

// SetStore is the setup function for a plugin
func (h *Simple) SetStore(s PluginStore) {
	h.db = s
}

// Descr describes this plugin
func (h *Simple) Descr() string {
	return "Return a new simple-login alias that can be used for various things."
}

// Re is what our simple-login matches
func (h *Simple) Re() string {
	return `(?i)^sl: (.+)$`
}

// Match checks for "simple-login: " messages
func (h *Simple) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Simple) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

func (h *Simple) Process(from, post string) (string, func() string) {
	reqInfo := h.fix(post)
	if reqInfo != "" {
		userKey := fmt.Appendf(nil, "simple_login_api_%s", from)
		userFile := base32.StdEncoding.EncodeToString(userKey)
		log.Println(userFile)
		userAPIKey, err := h.db.Get(userFile)
		if err != nil {
			return fmt.Sprintf("sorry %s, looks like you can't make aliases!", from), RespStub
		}

		reqURL, err := url.Parse("https://app.simplelogin.io/api/alias/random/new/")
		if err != nil {
			return fmt.Sprintf("sorry %s, invalid URL: %s", from, err), RespStub
		}
		v := url.Values{}
		v.Add("hostname", reqInfo)
		reqURL.RawQuery = v.Encode()

		var resp = &SimpleResp{}
		var req = HTTPRequest{
			Timeout: 15 * time.Second,
			URL:     reqURL.String(),
			Method:  "POST",
			ResBody: resp,
			Headers: map[string]string{
				"Authentication": userAPIKey,
			},
		}
		err = req.DoJSON()
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't hit the simple-login API. %s", from, err), RespStub
		}
		return resp.Email, RespStub
	}

	return "shrug.", RespStub
}

// RespondText to looking up of simple-login lookup requests
func (h *Simple) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name Simple!
func (h *Simple) Name() string {
	return "Simple"
}
