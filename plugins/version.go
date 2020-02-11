package plugins

import (
	"fmt"
	"log"
	"regexp"
	"runtime"

	"github.com/matrix-org/gomatrix"
)

// Version responds to hi messages
type Version struct {
}

// Descr describes this plugin
func (v *Version) Descr() string {
	return "Show a bit of information about what we are."
}

// Re matches version
func (v *Version) Re() string {
	return `(?i)version$`
}

// Match checks for "version" anywhere. Might want to tighten this one down at
// some point
func (v *Version) Match(user, msg string) bool {
	re := regexp.MustCompile(v.Re())
	return re.MatchString(msg) && ToMe(user, msg)
}

func (v *Version) print(to string) string {
	return fmt.Sprintf("%s, I am written in Go, running on %s", to, runtime.GOOS)
}

// SetStore does nothing in here
func (v *Version) SetStore(s PluginStore) {}

// RespondText to version events
func (v *Version) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	s := NameRE.ReplaceAllString(ev.Sender, "$1")

	log.Printf("%s: responding to '%s'", v.Name(), ev.Sender)
	SendText(c, ev.RoomID, v.print(s))
}

// Name Version
func (v *Version) Name() string {
	return "Version"
}
