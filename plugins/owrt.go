package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/matrix-org/gomatrix"
)

type OWRTData struct {
	Columns []string `json:"columns"`
	Entries [][]any  `json:"entries"`
}

func fetchJson(dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get("https://openwrt.org/toh.json")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func loadJson(name string) (*OWRTData, error) {
	d := &OWRTData{}

	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err = fetchJson(name)
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func getRow(d []any, c []int) []string {
	parts := []string{}
	for _, cidx := range c {
		val := fmt.Sprintf("%v", d[cidx])
		if strings.Contains(val, "toh:") {
			val = fmt.Sprintf("https://openwrt.org/%s", strings.Replace(val, ":", "/", 2))
		}
		parts = append(parts, val)
	}
	return parts
}

// OWRT lets one query openwrt's device db for compatible devices
type OWRT struct {
}

// Descr describes this plugin
func (h *OWRT) Descr() string {
	return "Query OpenWRT's DB for a device."
}

// Re is the regex for matching beat messages.
func (h *OWRT) Re() string {
	return `(?i)^owrt: |^openwrt: `
}

// Match determines if we are asking for a beat
func (h *OWRT) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

func (h *OWRT) fix(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

// SetStore we don't need a store here
func (h *OWRT) SetStore(_ PluginStore) {}

// RespondText to beat request events
func (h *OWRT) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	resp, _ := h.Process("", "")
	return SendText(c, ev.RoomID, resp)
}

// Process does the heavy lifting of calculating .beat
func (h *OWRT) Process(from, msg string) (string, func() string) {
	var (
		colSet = []int{}
		cols   = []string{
			"model",
			"brand",
			"version",
			"supportedcurrentrel",
			"devicepage",
		}
		device = strings.ToLower(h.fix(msg))
	)

	d, err := loadJson("/tmp/toh.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range cols {
		colSet = append(colSet, slices.Index(d.Columns, name))
	}

	rowEntries := make(map[string][]any)
	for _, row := range d.Entries {
		for _, entry := range row {
			strEntry := strings.ToLower(fmt.Sprintf("%s", entry))
			strRow := fmt.Sprintf("%s", row)
			if strings.Contains(strEntry, device) {
				if _, ok := rowEntries[strRow]; !ok {
					rowEntries[strRow] = row
				}
			}
		}
	}

	resp := []string{}
	for _, e := range rowEntries {
		resp = append(resp, strings.Join(getRow(e, colSet), ", "))
	}
	return strings.Join(resp, "\n"), RespStub
}

// Name beat
func (h *OWRT) Name() string {
	return "OWRT"
}
