package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// HomesteadResp is the json returned from our api
type HomesteadResp struct {
	Status string `json:"status"`
	Data   struct {
		Resulttype string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Name     string `json:"name"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// Homestead is our plugin type
type Homestead struct {
	db PluginStore
}

// SetStore is the setup function for a plugin
func (h *Homestead) SetStore(s PluginStore) {
	h.db = s
}

func (h *Homestead) get(loc string) (*HomesteadResp, error) {
	u := "https://graph.tapenet.org/_pub"
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w = &HomesteadResp{}
	err = json.Unmarshal(body, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Descr describes this plugin
func (h *Homestead) Descr() string {
	return "Display weather information for the Homestead"
}

// Re is what our weather matches
func (h *Homestead) Re() string {
	return `(?i)^home:|^homestead:\s?(\w+)?$`
}

// Match checks for "home: name?" messages
func (h *Homestead) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Homestead) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

func (h *Homestead) Process(from, post string) string {
	weather := h.fix(post)
	var s []string
	wd, err := h.get(weather)
	if err != nil {
		return fmt.Sprintf("sorry %s, I can't connect to the homestead. %q", from, err)
	}

	for _, e := range wd.Data.Result {
		if temp, err := strconv.ParseFloat(e.Value[1].(string), 64); err == nil {
			s = append(s, fmt.Sprintf("%s: %.2fC (%.2fF)", e.Metric.Name, temp, (temp*1.8000)+32.00))
		} else {
			log.Fatal(err)
		}
	}

	return strings.Join(s, ", ")
}

// RespondText to looking up of weather lookup requests
func (h *Homestead) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendText(c, ev.RoomID, h.Process(ev.Sender, post))
}

// Name Homestead!
func (h *Homestead) Name() string {
	return "Homestead"
}
