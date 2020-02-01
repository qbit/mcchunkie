package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Beer responds to hi messages
type Beer struct {
}

// BeerRecord is the parent beer type
type BeerRecord struct {
	DatasetID string     `json:"datasetid"`
	RecordID  string     `json:"recordid"`
	Fields    BeerFields `json:"fields"`
}

// BeerFields ar the fields we care about for any given beer
type BeerFields struct {
	Website     string `json:"website"`
	City        string `json:"city"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Descr       string `json:"descript"`
	Style       string `json:"style_name"`
	IBU         int    `json:"ibu"`
	BreweryName string `json:"name_breweries"`
}

// BeerRecords are a collection of responses
type BeerRecords []BeerRecord

func (h *Beer) fix(msg string) string {
	re := regexp.MustCompile(`(?i)^beer: `)
	return re.ReplaceAllString(msg, "$1")
}

func (h *Beer) match(msg string) bool {
	re := regexp.MustCompile(`(?i)^beer: `)
	return re.MatchString(msg)
}

func (h *Beer) get(beer string) (*BeerRecords, error) {
	u := "https://data.opendatasoft.com/api/records/1.0/search?dataset=open-beer-database%40public-us&q="
	u = fmt.Sprintf("%s%s", u, url.PathEscape(beer))
	log.Println(u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	beers := &BeerRecords{}
	_ = json.Unmarshal([]byte(body), beers)

	return beers, nil
}

// Respond to hi events
func (h *Beer) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				if h.match(post) {
					beer := h.fix(post)
					if beer != "" {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						j, err := h.get(beer)
						if err != nil {
							SendMessage(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look for beer. (%s)", ev.Sender, err))
						}
						SendMessage(c, ev.RoomID, fmt.Sprintf("%v", j))
					}
				}
			}
		}
	}
}

// Name hi
func (h *Beer) Name() string {
	return "Beer"
}
