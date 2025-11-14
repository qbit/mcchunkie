package plugins

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Beer responds to beer requests
type Beer struct {
	store PluginStore
}

type BeerResps struct {
	Code  int        `json:"code"`
	Error bool       `json:"error"`
	Data  []BeerData `json:"data"`
}

type BeerResp struct {
	Code  int      `json:"code"`
	Error bool     `json:"error"`
	Data  BeerData `json:"data"`
}

type BeerData struct {
	Sku                    string `json:"sku"`
	Name                   string `json:"name"`
	Brewery                string `json:"brewery"`
	Rating                 string `json:"rating"`
	Category               string `json:"category"`
	SubCategory1           string `json:"sub_category_1"`
	SubCategory2           string `json:"sub_category_2"`
	SubCategory3           string `json:"sub_category_3"`
	Description            string `json:"description"`
	Region                 string `json:"region"`
	Country                string `json:"country"`
	Abv                    string `json:"abv"`
	Ibu                    string `json:"ibu"`
	CaloriesPerServing12Oz string `json:"calories_per_serving_12oz"`
	CarbsPerServing12Oz    string `json:"carbs_per_serving_12oz"`
	TastingNotes           string `json:"tasting_notes"`
	FoodPairing            string `json:"food_pairing"`
	SuggestedGlassware     string `json:"suggested_glassware"`
	ServingTempF           string `json:"serving_temp_f"`
	ServingTempC           string `json:"serving_temp_c"`
	BeerType               string `json:"beer_type"`
	Features               string `json:"features"`
}

// Descr describes this plugin
func (h *Beer) Descr() string {
	return "Queries [Winevybe](https://winevybe.com/beer-api)'s beer database for a given beer."
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
func (h *Beer) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Beer) pretty(d BeerData) string {
	if d.Description != "" {
		d.Description = "\n" + d.Description
	}
	return fmt.Sprintf("%s (%s) by %s from %s, %s - IBU: %s, ABV: %s %s",
		d.Name,
		d.SubCategory3,
		d.Brewery,
		d.Region,
		d.Country,
		d.Ibu,
		d.Abv,
		d.Description,
	)
}

// SetStore we don't need a store here.
func (h *Beer) SetStore(s PluginStore) {
	h.store = s
}

func (h *Beer) Process(from, msg string) (string, func() string) {
	key, _ := h.store.Get("beer_api_key")
	beer := h.fix(msg)
	resp := "¯\\_(ツ)_/¯"
	if beer != "" {
		u := fmt.Sprintf("%s%s",
			"https://beer9.p.rapidapi.com?name=",
			url.PathEscape(beer),
		)
		req := HTTPRequest{
			Method: "GET",
			URL:    u,
			Headers: map[string]string{
				"x-rapidapi-key":  key,
				"x-rapidapi-host": "beer9.p.rapidapi.com",
			},
		}

		data, err := req.Do()
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look for beer. (%s)", from, err), RespStub
		}

		var singleBeer BeerResp
		var multipleBeer BeerResps
		err = json.Unmarshal(data, &multipleBeer)
		if err == nil && len(multipleBeer.Data) > 0 {
			l := len(multipleBeer.Data)
			rb := multipleBeer.Data[rand.Intn(l)]
			return fmt.Sprintf("Found %d results, here's a random one: %s", l, h.pretty(rb)), RespStub
		}
		err = json.Unmarshal(data, &singleBeer)
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look for beer. (%s)", from, err), RespStub
		}

		if singleBeer.Code == 200 {
			return h.pretty(singleBeer.Data), RespStub
		}

		return fmt.Sprintf("Sorry that beer is %d", singleBeer.Code), RespStub
	}
	return resp, RespStub
}

// RespondText to looking up of beer requests
func (h *Beer) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name Beer!
func (h *Beer) Name() string {
	return "Beer"
}
