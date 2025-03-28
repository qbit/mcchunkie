package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
	var s []string
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

type PollutionResp struct {
	Coord CoordResp `json:"coord"`
	List  []struct {
		Dt   int `json:"dt"`
		Main struct {
			Aqi int `json:"aqi"`
		} `json:"main"`
		Components struct {
			Co   float64 `json:"co"`
			No   float64 `json:"no"`
			No2  float64 `json:"no2"`
			O3   float64 `json:"o3"`
			So2  float64 `json:"so2"`
			Pm25 float64 `json:"pm2_5"`
			Pm10 float64 `json:"pm10"`
			Nh3  float64 `json:"nh3"`
		} `json:"components"`
	} `json:"list"`
}

func (p *PollutionResp) String() string {
	if len(p.List) == 0 {
		return "AQI: unavailable"
	}

	c := p.List[0].Components
	return fmt.Sprintf("AQI: %d (CO2: %.1f, PM2.5: %.1f, PM10: %.1f)", p.List[0].Main.Aqi, c.Co, c.Pm25, c.Pm10)
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

func (h *Weather) getPollution(c *CoordResp) (*PollutionResp, error) {
	u, err := url.Parse("http://api.openweathermap.org/data/2.5/air_pollution")
	if err != nil {
		return nil, err
	}
	key, err := h.db.Get("weather_api_key")
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, fmt.Errorf("no API key set")
	}

	v := url.Values{}
	v.Set("APPID", key)
	v.Add("lat", strconv.FormatFloat(c.Lat, 'g', -1, 64))
	v.Add("lon", strconv.FormatFloat(c.Lon, 'g', -1, 64))

	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w = &PollutionResp{}
	err = json.Unmarshal(body, w)
	if err != nil {
		log.Println(string(body))
		return nil, err
	}

	return w, nil
}

func (h *Weather) getCurrent(loc string) (*WeatherResp, error) {
	u := "http://api.openweathermap.org/data/2.5/weather?%s"
	key, err := h.db.Get("weather_api_key")
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, fmt.Errorf("no API key set")
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w = &WeatherResp{}
	err = json.Unmarshal(body, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Descr describes this plugin
func (h *Weather) Descr() string {
	return "Produce weather information for a given ZIP or ZIP,CountryCode combo. Data comes from openweathermap.org."
}

// Re is what our weather matches
func (h *Weather) Re() string {
	return `(?i)^weather: (\d+|\d+(:?,[a-z][A-Z]))$`
}

// Match checks for "weather: " messages
func (h *Weather) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *Weather) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

func (h *Weather) Process(from, post string) (string, func() string) {
	weather := h.fix(post)
	if weather != "" {
		wd, err := h.getCurrent(weather)
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look up the weather. %s", from, err), RespStub
		}
		po, err := h.getPollution(&wd.Coord)
		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look up the pollution. %s", from, err), RespStub
		}

		pollution := po.String()

		return fmt.Sprintf(`%s: %s (%s) %s, Humidity: %s%%, %s`,
			wd.Name,
			wd.c(),
			wd.f(),
			wd.conditions(),
			wd.humidity(),
			pollution,
		), RespStub
	}

	return "shrug.", RespStub
}

// RespondText to looking up of weather lookup requests
func (h *Weather) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, _ := h.Process(ev.Sender, post)
	return SendText(c, ev.RoomID, resp)
}

// Name Weather!
func (h *Weather) Name() string {
	return "Weather"
}
