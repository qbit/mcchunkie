package plugins

import (
	"fmt"
	"math/rand"
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

// Descr describes this plugin
func (h *Beer) Descr() string {
	return "Queries [OpenDataSoft](https://public-us.opendatasoft.com/explore/dataset/open-beer-database/table/)'s beer database for a given beer."
}

// Re returns the beer matching string
func (h *Beer) Re() string {
	return `(?i)^beer: `
}

func (h *Beer) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// Match determines if we should call the response for Beer
func (h *Beer) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
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

// SetStore we don't need a store here.
func (h *Beer) SetStore(s PluginStore) {}

// RespondText to looking up of beer requests
func (h *Beer) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	beer := h.fix(post)
	if beer != "" {
		var beers = &BeerResp{}
		u := fmt.Sprintf("%s%s",
			"https://data.opendatasoft.com/api/records/1.0/search?dataset=open-beer-database%40public-us&q=",
			url.PathEscape(beer),
		)
		req := HTTPRequest{
			Method:  "GET",
			ResBody: beers,
			URL:     u,
		}
		err := req.DoJSON()
		if err != nil {
			SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look for beer. (%s)", ev.Sender, err))
		}

		switch {
		case beers.Nhits == 0:
			SendText(c, ev.RoomID, "¯\\_(ツ)_/¯")
		case beers.Nhits == 1:
			SendText(c, ev.RoomID, h.pretty(*beers, false))
		case beers.Nhits > 1:
			SendText(c, ev.RoomID, fmt.Sprintf("Found %d beers, here is a random one:\n%s", beers.Nhits, h.pretty(*beers, true)))
		}
	}
}

// Name Beer!
func (h *Beer) Name() string {
	return "Beer"
}
