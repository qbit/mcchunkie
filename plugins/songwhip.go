package plugins

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

type SongwhipReq struct {
	URL string `json:"url"`
}

type SongwhipResp struct {
	Type     string `json:"type"`
	ID       int    `json:"id"`
	Path     string `json:"path"`
	PagePath string `json:"pagePath"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Links    struct {
		Qobuz        bool `json:"qobuz"`
		Tidal        bool `json:"tidal"`
		Amazon       bool `json:"amazon"`
		Deezer       bool `json:"deezer"`
		Itunes       bool `json:"itunes"`
		Discogs      bool `json:"discogs"`
		Napster      bool `json:"napster"`
		Pandora      bool `json:"pandora"`
		Spotify      bool `json:"spotify"`
		Twitter      bool `json:"twitter"`
		Youtube      bool `json:"youtube"`
		Bandcamp     bool `json:"bandcamp"`
		Facebook     bool `json:"facebook"`
		Audiomack    bool `json:"audiomack"`
		Instagram    bool `json:"instagram"`
		LineMusic    bool `json:"lineMusic"`
		Soundcloud   bool `json:"soundcloud"`
		AmazonMusic  bool `json:"amazonMusic"`
		ItunesStore  bool `json:"itunesStore"`
		MusicBrainz  bool `json:"musicBrainz"`
		YoutubeMusic bool `json:"youtubeMusic"`
	} `json:"links"`
	Description          interface{} `json:"description"`
	LinksCountries       []string    `json:"linksCountries"`
	SourceCountry        string      `json:"sourceCountry"`
	SpotifyID            string      `json:"spotifyId"`
	CreatedAtTimestamp   int64       `json:"createdAtTimestamp"`
	RefreshedAtTimestamp int64       `json:"refreshedAtTimestamp"`
	URL                  string      `json:"url"`
}

type Songwhip struct{}

// Descr describes this plugin
func (s *Songwhip) Descr() string {
	return "Get a songwhip link for a music link"
}

// Re returns the federation check matching string
func (s *Songwhip) Re() string {
	return `(?i)^(?:songwhip: |sw: |music: )(.*)$`
}

func (s *Songwhip) fix(msg string) string {
	re := regexp.MustCompile(s.Re())
	return re.ReplaceAllString(msg, "$1")
}

// Match determines if we should call the response for Songwhip
func (s *Songwhip) Match(_, msg string) bool {
	re := regexp.MustCompile(s.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (s *Songwhip) SetStore(_ PluginStore) {}

func hasService(s string, enabled bool) string {
	if enabled {
		return fmt.Sprintf("%s ✅", s)
	}

	return fmt.Sprintf("%s 🇽", s)
}

func (s *Songwhip) Process(from, post string) string {
	musicURL := s.fix(post)
	if musicURL != "" {
		_, err := url.Parse(musicURL)
		if err != nil {
			return fmt.Sprintf("Please don't abuse this free service. that's not a real url: %q", musicURL)
		}

		var swresp = &SongwhipResp{}

		var req = HTTPRequest{
			Timeout: 5 * time.Second,
			URL:     "https://songwhip.com/",
			Method:  "POST",
			ResBody: swresp,
			ReqBody: &SongwhipReq{
				URL: musicURL,
			},
		}
		err = req.DoJSON()

		if err != nil {
			return fmt.Sprintf("sorry %s, I can't look up that link on songwhip (%q)", from, err)
		}

		return fmt.Sprintf("[%s](%s) (%s) can be found on: %s, %s, %s, %s",
			swresp.Name,
			swresp.URL,
			swresp.Type,
			hasService("AppleMusic", swresp.Links.Itunes),
			hasService("Spotify", swresp.Links.Spotify),
			hasService("Tidal", swresp.Links.Tidal),
			hasService("YTMusic", swresp.Links.YoutubeMusic),
		)
	}
	return "invalid hostname"
}

// RespondText to looking up of federation check requests
func (s *Songwhip) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendMD(c, ev.RoomID, s.Process(ev.Sender, post))
}

// Name Songwhip!
func (s *Songwhip) Name() string {
	return "Songwhip"
}
