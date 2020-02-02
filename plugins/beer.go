package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// Beer responds to beer requests
type Beer struct {
}

// BeerResp represents our response from the api
type BeerResp struct {
	Nhits      int        `json:"nhits"`
	Parameters Parameters `json:"parameters"`
	Records    []Records  `json:"records"`
}

// Parameters are the meta information
type Parameters struct {
	Dataset  []string `json:"dataset"`
	Timezone string   `json:"timezone"`
	Q        string   `json:"q"`
	Rows     int      `json:"rows"`
	Format   string   `json:"format"`
}

// Fields are the bits of info we care about
type Fields struct {
	Website       string    `json:"website"`
	City          string    `json:"city"`
	StyleID       string    `json:"style_id"`
	Name          string    `json:"name"`
	Country       string    `json:"country"`
	CatID         string    `json:"cat_id"`
	BreweryID     string    `json:"brewery_id"`
	Descript      string    `json:"descript"`
	Upc           int       `json:"upc"`
	Coordinates   []float64 `json:"coordinates"`
	Ibu           int       `json:"ibu"`
	CatName       string    `json:"cat_name"`
	LastMod       time.Time `json:"last_mod"`
	State         string    `json:"state"`
	StyleName     string    `json:"style_name"`
	Abv           float64   `json:"abv"`
	Address1      string    `json:"address1"`
	NameBreweries string    `json:"name_breweries"`
	Srm           int       `json:"srm"`
	ID            string    `json:"id"`
	AddUser       string    `json:"add_user"`
}

// Geometry is basically useless
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Records holds our fileds
type Records struct {
	Datasetid       string    `json:"datasetid"`
	Recordid        string    `json:"recordid"`
	Fields          Fields    `json:"fields"`
	Geometry        Geometry  `json:"geometry"`
	RecordTimestamp time.Time `json:"record_timestamp"`
}

func (h *Beer) fix(msg string) string {
	re := regexp.MustCompile(`(?i)^beer: `)
	return re.ReplaceAllString(msg, "$1")
}

func (h *Beer) match(msg string) bool {
	re := regexp.MustCompile(`(?i)^beer: `)
	return re.MatchString(msg)
}

func (h *Beer) get(beer string) (*BeerResp, error) {
	u := "https://data.opendatasoft.com/api/records/1.0/search?dataset=open-beer-database%40public-us&q="
	u = fmt.Sprintf("%s%s", u, url.PathEscape(beer))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var beers = &BeerResp{}
	err = json.Unmarshal([]byte(body), beers)
	if err != nil {
		return nil, err
	}

	return beers, nil
}

func (h *Beer) pretty(b BeerResp, random bool) string {
	idx := 0

	if random {
		rand.Seed(time.Now().Unix())
		idx = rand.Intn(len(b.Records))
	}

	return fmt.Sprintf("%s (%s) by %s from %s, %s - IBU: %d, ABV: %.1f %s\n%s",
		b.Records[idx].Fields.Name,
		b.Records[idx].Fields.StyleName,
		b.Records[idx].Fields.NameBreweries,
		b.Records[idx].Fields.City,
		b.Records[idx].Fields.State,
		b.Records[idx].Fields.Ibu,
		b.Records[idx].Fields.Abv,
		b.Records[idx].Fields.Website,
		b.Records[idx].Fields.Descript,
	)
}

// Respond to looking up of beer requests
func (h *Beer) Respond(c *gomatrix.Client, ev *gomatrix.Event, user string) {
	if mtype, ok := ev.MessageType(); ok {
		switch mtype {
		case "m.text":
			if post, ok := ev.Body(); ok {
				if h.match(post) {
					beer := h.fix(post)
					if beer != "" {
						log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
						brr, err := h.get(beer)
						if err != nil {
							SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look for beer. (%s)", ev.Sender, err))
						}

						switch {
						case brr.Nhits == 0:
							SendText(c, ev.RoomID, "¯\\_(ツ)_/¯")
						case brr.Nhits == 1:
							SendText(c, ev.RoomID, fmt.Sprintf("%s", h.pretty(*brr, false)))
						case brr.Nhits > 1:
							SendText(c, ev.RoomID, fmt.Sprintf("Found %d beers, here is a random one:\n%s", brr.Nhits, h.pretty(*brr, true)))
						}
					}
				}
			}
		}
	}
}

// Name Beer!
func (h *Beer) Name() string {
	return "Beer"
}
