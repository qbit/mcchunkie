package plugins

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
)

// PluginStore matches MCStore. This allows the main store to be used by
// plugins.
type PluginStore interface {
	Set(key, values string)
	Get(key string) (string, error)
}

// Plugin defines the interface a plugin must implement to be used by
// mcchunkie.
type Plugin interface {
	// Descr returns a brief description of the plugin.
	Descr() string

	// Match determines if the plugin's main Respond function should be
	// called
	Match(user, message string) bool

	// Name should return the human readable name of the bot
	Name() string

	// Re returns the regular expression that a plugin uses to "match"
	Re() string

	// RespondText responds to a "m.text" event
	RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, path string)

	// SetStore exposes the top level MCStore to a plugin
	SetStore(s PluginStore)
}

// NameRE matches the "friendly" name. This is typically used in tab
// completion.
var NameRE = regexp.MustCompile(`@(.+):.+$`)

// ToMe returns true of the message pertains to the bot
func ToMe(user, message string) bool {
	u := NameRE.ReplaceAllString(user, "$1")
	return strings.Contains(message, u)
}

// HTTPRequest has the bits for making http requests
type HTTPRequest struct {
	Client  http.Client
	Request *http.Request
	Timeout time.Duration
	URL     string
	Method  string
	ReqBody interface{}
	ResBody interface{}
}

// DoJSON is a general purpose http mechanic that can be used to get, post..
// what evs. The response is always expected to be json
func (h *HTTPRequest) DoJSON() (err error) {
	h.Client.Timeout = h.Timeout

	if h.Method == "" {
		h.Method = "GET"
	}

	if h.ReqBody != nil {
		// We have a request to send to the server
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(h.ReqBody)
		if err != nil {
			return err
		}
		h.Request, err = http.NewRequest(h.Method, h.URL, buf)
	} else {
		// Just gimme dem datas
		h.Request, err = http.NewRequest(h.Method, h.URL, nil)
	}

	if err != nil {
		return err
	}

	h.Request.Header.Set("Content-Type", "application/json")

	res, err := h.Client.Do(h.Request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}

	if h.ResBody != nil && res.Body != nil {
		return json.NewDecoder(res.Body).Decode(&h.ResBody)
	}

	return nil
}

// SendText sends a text message to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendText(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	c.SendText(roomID, message)

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// Plugins is a collection of our plugins. An instance of this is iterated
// over for each message the bot receives.
type Plugins []Plugin

// Plugs defines the "enabled" plugins.
var Plugs = Plugins{
	&Beat{},
	&Beer{},
	&BotSnack{},
	&Feder{},
	&HighFive{},
	&Hi{},
	&LoveYou{},
	&OpenBSDMan{},
	&Source{},
	&Thanks{},
	&Version{},
	&Wb{},
	&Weather{},
}
