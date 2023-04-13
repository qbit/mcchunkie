package plugins

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
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

	// RespondMatrix responds to a Matrix "m.text" event
	RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, path string) error

	// Process is the processed response from the plugin. This is useful for
	// running the plugins outside of the context of Matrix.
	Process(from, message string) string

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

// RemoveName removes the friendly name from a given message
func RemoveName(user, message string) string {
	n := NameRE.ReplaceAllString(user, "$1")
	return strings.ReplaceAll(message, n+": ", "")
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

	_, err = c.SendText(roomID, message)
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendEmote sends an emote to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendEmote(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	_, err = c.SendMessageEvent(roomID, "m.room.message", gomatrix.GetHTMLMessage("m.emote", message))
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendHTML sends an html message to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendHTML(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	_, err = c.SendMessageEvent(roomID, "m.room.message", gomatrix.GetHTMLMessage("m.text", message))
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendMDNotice sends an html notice to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendMDNotice(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	md := []byte(message)
	html := markdown.ToHTML(md, nil, nil)
	_, err = c.SendMessageEvent(roomID, "m.room.message", gomatrix.GetHTMLMessage("m.notice", string(html)))
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendUnescNotice sends an text notice to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendUnescNotice(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	// Undo the escaping
	_, err = c.SendMessageEvent(roomID, "m.room.message", gomatrix.HTMLMessage{
		Body:          message,
		MsgType:       "m.notice",
		Format:        "org.matrix.custom.text",
		FormattedBody: message,
	})
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendNotice sends an text notice to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendNotice(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	_, err = c.SendMessageEvent(roomID, "m.room.message", gomatrix.GetHTMLMessage("m.notice", message))
	if err != nil {
		return err
	}

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// SendMD takes markdown and converts it to an html message.
func SendMD(c *gomatrix.Client, roomID, message string) error {
	md := []byte(message)
	html := markdown.ToHTML(md, nil, nil)
	return SendHTML(c, roomID, string(html))
}

// SendImage takes an image and sends it!.
func SendImage(c *gomatrix.Client, roomID string, img *image.RGBA) error {
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		_ = png.Encode(w, img)
	}()

	mediaURL, err := c.UploadToContentRepo(r, "image/png", 0)
	if err != nil {
		return err
	}

	_, err = c.SendImage(roomID, "embedded_image.png", mediaURL.ContentURI)
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
	&Ban{},
	&BananaStab{},
	&Beat{},
	&Beer{},
	&BotSnack{},
	&DMR{},
	&Feder{},
	&Groan{},
	&Ham{},
	&Hi{},
	&HighFive{},
	&Homestead{},
	&LoveYou{},
	&OpenBSDMan{},
	&Palette{},
	&PGP{},
	&RFC{},
	&Salute{},
	&Snap{},
	&Songwhip{},
	&Source{},
	&Thanks{},
	&Toki{},
	&Version{},
	&Wb{},
	&Weather{},
	&Yeah{},
}
