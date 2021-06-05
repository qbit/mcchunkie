package chats

import (
	"suah.dev/mcchunkie/plugins"
)

// Message is a cross-service representation of a chat message
type Message struct {
	Service string // TODO: Maybe this should be a type?
	To      string
	From    string
	Body    string
}

// ChatStore matches MCStore. This allows the main store to be used by
// plugins.
type ChatStore interface {
	Set(key, values string)
	Get(key string) (string, error)
}

// Chat represents a mode of communication like Matrix, IRC or SMS.
type Chat interface {
	Connect(s plugins.PluginStore) error
	Process(m *Message) error
	Disconnect() error
}

// Chats is a collection of our chat methods. An instance of this is iterated
// over for each message the bot responds to.
type Chats []Chat

// ChatMethods defines the "enabled" chat methogs.
var ChatMethods = Chats{
	&Matrix{},
	&IRC{},
}
