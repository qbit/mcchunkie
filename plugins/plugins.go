package plugins

import (
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// Plugin is an interface that specifies what a plugin needs to respond to.
type Plugin interface {
	Respond(c *gomatrix.Client, ev *gomatrix.Event, user string)
	Name() string
}

// NameRE matches just the name of a matrix user
var NameRE = regexp.MustCompile(`@(.+):.+$`)

// ToMe returns true of the message pertains to the bot
func ToMe(user, message string) bool {
	return strings.Contains(message, user)
}

// SendMessage sends a message to a given room
func SendMessage(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	c.SendText(roomID, message)

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

// Plugins area  collection
type Plugins []Plugin

// Plugs are all of our plugins
var Plugs = Plugins{
	&HighFive{},
	&Hi{},
}
