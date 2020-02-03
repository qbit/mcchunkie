package plugins

import (
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// Plugin defines the functions a plugin must implement to be used by
// mcchunkie.
type Plugin interface {
	//Respond(c *gomatrix.Client, ev *gomatrix.Event, user string)
	RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, path string)
	Name() string
}

// NameRE matches the "friendly" name. This is typically used in tab
// completion.
var NameRE = regexp.MustCompile(`@(.+):.+$`)

// ToMe returns true of the message pertains to the bot
func ToMe(user, message string) bool {
	return strings.Contains(message, user)
}

// SendText sends a text message to a given room. It pretends to be
// "typing" by calling UserTyping for the caller.
func SendText(c *gomatrix.Client, roomID, message string) error {
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

// Plugins is a collection of our plugins. An instance of this is iterated
// over for each message the bot receives.
type Plugins []Plugin

// Plugs defines the "enabled" plugins.
var Plugs = Plugins{
	&Beer{},
	&BotSnack{},
	&HighFive{},
	&Hi{},
	&LoveYou{},
	&OpenBSDMan{},
	&Source{},
	&Version{},
}
