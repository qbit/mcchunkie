package plugins

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
	"golang.org/x/crypto/openpgp"
)

// PGP is our plugin type
type PGP struct {
}

// SetStore is the setup function for a plugin
func (p *PGP) SetStore(s PluginStore) {
}

// Descr describes this plugin
func (p *PGP) Descr() string {
	return "Queries keys.openpgp.org"
}

// Re is what our pgp request matches
func (p *PGP) Re() string {
	return `(?i)^pgp: (.+@.+\..+|[a-f0-9]+)$`
}

// Match checks for "pgp: " messages
func (p *PGP) Match(_, msg string) bool {
	re := regexp.MustCompile(p.Re())
	return re.MatchString(msg)
}

func (p *PGP) fix(msg string) string {
	re := regexp.MustCompile(p.Re())
	return strings.ToUpper(re.ReplaceAllString(msg, "$1"))
}

// RespondText to looking up of PGP info
func (p *PGP) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	search := p.fix(post)
	searchURL := "https://keys.openpgp.org//vks/v1/by-fingerprint/%s"

	if strings.ContainsAny(search, "@") {
		searchURL = "https://keys.openpgp.org//vks/v1/by-email/%s"
	}

	escSearch, err := url.Parse(search)
	if err != nil {
		return err
	}

	u := fmt.Sprintf(searchURL, escSearch)

	resp, err := http.Get(u)
	if err != nil {
		_ = SendText(c, ev.RoomID, fmt.Sprintf("Can't search keys.openpgp.org: %s", err))
		return err
	}

	defer resp.Body.Close()

	kr, err := openpgp.ReadArmoredKeyRing(resp.Body)
	if err != nil {
		_ = SendText(c, ev.RoomID, fmt.Sprintf("Can't parse key ring: %s", err))
		return err
	}

	var ids []string
	var fps []string
	for _, entity := range kr {
		for _, i := range entity.Identities {
			ids = append(ids, fmt.Sprintf("- %q", i.Name))
		}
		fps = append(fps, fmt.Sprintf("**Fingerprint**: %s",
			hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])))
	}

	return SendMD(c, ev.RoomID, fmt.Sprintf("%s\n\n%s",
		strings.Join(ids, "\n"),
		strings.Join(fps, "\n")))
}

// Name PGP!
func (p *PGP) Name() string {
	return "PGP"
}
