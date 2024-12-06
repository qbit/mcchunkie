package chats

import (
	"strings"

	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
)

// Chat represents a mode of communication like Matrix, IRC or SMS.
type Chat interface {
	// Connect connects
	Connect(*mcstore.MCStore, *plugins.Plugins) error
	Name() string
	Send(string, string) error
}

// Chats is a collection of our chat methods. An instance of this is iterated
// over for each message the bot responds to.
type Chats []Chat

// ChatMethods defines the "enabled" chat methogs.
var ChatMethods = Chats{
	&MatrixChat{},
	&XMPPChat{},
	&IRCChat{},
	&MailChat{},
	&SMSChat{},
}

func (c *Chats) List() string {
	s := []string{}

	for _, ch := range *c {
		s = append(s, ch.Name())
	}

	return strings.Join(s, ", ")
}
