package plugins

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/matrix-org/gomatrix"
)

var version string

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

func (v *Version) print() string {
	if version == "" {
		version = "unknown version"
	}
	return fmt.Sprintf("%s running on %s", version, runtime.GOOS)
}

// SetStore does nothing in here
func (v *Version) SetStore(_ PluginStore) {}

// RespondText to version events
func (v *Version) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	return SendText(c, ev.RoomID, v.print())
}

// Name Version
func (v *Version) Name() string {
	return "Version"
}
