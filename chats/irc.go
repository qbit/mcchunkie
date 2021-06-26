package chats

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"gopkg.in/irc.v3"
	"suah.dev/mcchunkie/plugins"
)

// IRC represents our IRC chat service
type IRC struct {
	Client irc.Client
	Store  plugins.PluginStore
	Log    *log.Logger
}

// Connect connects to an IRC server
func (i *IRC) Connect(store plugins.PluginStore) error {
	var ircServer, _ = store.Get("irc_server")
	var ircPort, _ = store.Get("irc_port")
	var ircNick, _ = store.Get("irc_nick")
	var ircPass, _ = store.Get("irc_pass")
	var ircRooms, _ = store.Get("irc_rooms")
	//var toRE = regexp.MustCompile(`^:(\w+)\s`)

	i.Log = log.Default()
	i.Log.SetPrefix("IRC: ")

	i.Store = store

	if ircServer != "" {

		i.Log.Printf("connecting to %q\n", ircServer)

		dialStr := fmt.Sprintf("%s:%s", ircServer, ircPort)
		conn, err := tls.Dial("tcp", dialStr, &tls.Config{
			ServerName: ircServer,
		})
		if err != nil {
			return err
		}

		config := irc.ClientConfig{
			Nick: ircNick,
			Pass: ircPass,
			User: ircNick,
			Name: "McChunkie",
			Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
				switch m.Command {
				case "001":
					for _, r := range strings.Split(ircRooms, ",") {
						i.Log.Printf("joining %q\n", r)
						c.Write(fmt.Sprintf("JOIN %s", r))
					}
				case "PING":
					server := m.Trailing()
					i.Log.Printf("pong %q\n", server)
					c.Write(fmt.Sprintf("PONG %s", server))
				case "INVITE":
					room := m.Trailing()
					i.Log.Printf("joining %q\n", room)
					c.Write(fmt.Sprintf("JOIN %s", room))
				case "PRIVMSG":
					msg := m.Trailing()
					from := m.Prefix.Name
					to := m.Params[0]

					if from == c.CurrentNick() {
						// Ignore messages from ourselves
						return
					}

					resp := ""
					for _, p := range plugins.Plugs {
						if p.Match(c.CurrentNick(), msg) {
							p.SetStore(store)

							resp = p.Process(from, msg)
						}
					}

					if !c.FromChannel(m) {
						// in a private chat
						to = from

					}

					if resp != "" {
						i.Log.Printf("sending: %q to %q\n", resp, to)
						c.WriteMessage(&irc.Message{
							Command: "PRIVMSG",
							Params: []string{
								to,
								resp,
							},
						})
					}
				default:
					i.Log.Printf("unhandled - %q", m.String())
				}
			}),
		}

		client := irc.NewClient(conn, config)
		err = client.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// Disconnect disconnects cleanely from our service
func (i *IRC) Disconnect() error {
	return nil
}

// Process receives a Message and determines if it should be used or not
func (i *IRC) Process(msg *Message) error {
	if msg.Service != "IRC" {
		return nil
	}
	return nil
}
