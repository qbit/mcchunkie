package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// WeatherResp is a JSON response from OpenWeatherMap.org
type WeatherResp struct {
	Coord      CoordResp   `json:"coord"`
	Weather    []FeelsResp `json:"weather"`
	Base       string      `json:"base"`
	Main       MainResp    `json:"main"`
	Visibility int         `json:"visibility"`
	Wind       WindResp    `json:"wind"`
	Clouds     CloudsResp  `json:"clouds"`
	Dt         int         `json:"dt"`
	Sys        SysResp     `json:"sys"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Cod        int         `json:"cod"`
}

func (w *WeatherResp) conditions() string {
	s := []string{}
	for _, cond := range w.Weather {
		s = append(s, cond.Description)
	}
	return strings.Join(s[:], ", ")
}

func (w *WeatherResp) f() string {
	return fmt.Sprintf("%.1f °F", (w.Main.Temp-273.15)*1.8000+32.00)
}

func (w *WeatherResp) c() string {
	return fmt.Sprintf("%.1f °C", w.Main.Temp-273.15)
}

func (w *WeatherResp) humidity() string {
	return fmt.Sprintf("%d", w.Main.Humidity)
}

// CoordResp is the log/lat of the location queried
type CoordResp struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// FeelsResp represents the friendly info like "light rain"
type FeelsResp struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// MainResp contains the bulk of the weather data
type MainResp struct {
	Temp     float64 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

// WindResp gives us various bits of information about the wind. Direction, etc.
type WindResp struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// CloudsResp ?
type CloudsResp struct {
	All int `json:"all"`
}

// SysResp seems to be used internally for something
type SysResp struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// Weather is our plugin type
type Weather struct {
	db PluginStore
}

// SetStore is the setup function for a plugin
func (h *Weather) SetStore(s PluginStore) {
	h.db = s
}

func (h *Weather) get(loc string) (*WeatherResp, error) {
	u := "http://api.openweathermap.org/data/2.5/weather?%s"
	key, err := h.db.Get("weather_api_key")
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, fmt.Errorf("No API key set")
	}

	v := url.Values{}
	v.Set("APPID", key)
	v.Add("zip", loc)

	u = fmt.Sprintf(u, v.Encode())

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w = &WeatherResp{}
	err = json.Unmarshal([]byte(body), w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (h *Weather) match(msg string) bool {
	re := regexp.MustCompile(`(?i)^weather: \d+$`)
	return re.MatchString(msg)
}

func (h *Weather) fix(msg string) string {
	re := regexp.MustCompile(`(?i)^weather: (\d+)$`)
	return re.ReplaceAllString(msg, "$1")
}

// RespondText to looking up of weather lookup requests
func (h *Weather) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	if h.match(post) {
		weather := h.fix(post)
		if weather != "" {
			log.Printf("%s: responding to '%s'", h.Name(), ev.Sender)
			wd, err := h.get(weather)
			if err != nil {
				SendText(c, ev.RoomID, fmt.Sprintf("sorry %s, I can't look up the weather. %s", ev.Sender, err))
			}
			SendText(c, ev.RoomID,
				fmt.Sprintf("%s: %s (%s) Humidity: %s, %s",
					wd.Name,
					wd.f(),
					wd.c(),
					wd.humidity(),
					wd.conditions(),
				))
		}
	}
}

// Name Weather!
func (h *Weather) Name() string {
	return "Weather"
}
