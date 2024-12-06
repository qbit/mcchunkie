package chats

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"gopkg.in/irc.v3"
	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
)

type IRCChat struct{}

func (i *IRCChat) Name() string {
	return "IRC"
}

func (i *IRCChat) Send(string, string) error {
	return nil
}

// IRCConnect connects to our irc server
func (i *IRCChat) Connect(store *mcstore.MCStore, plugins *plugins.Plugins) error {
	ircServer, err := store.Get("irc_server")
	if err != nil {
		return err
	}
	ircPort, err := store.Get("irc_port")
	if err != nil {
		return err
	}
	ircNick, err := store.Get("irc_nick")
	if err != nil {
		return err
	}
	ircPass, err := store.Get("irc_pass")
	if err != nil {
		log.Println(err)
	}
	ircRooms, err := store.Get("irc_rooms")
	if err != nil {
		return err
	}

	if ircServer != "" {

		log.Printf("IRC: connecting to %q\n", ircServer)

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
						log.Printf("IRC: joining %q\n", r)
						c.Write(fmt.Sprintf("JOIN %s", r))
					}
				case "PING":
					server := m.Trailing()
					log.Printf("IRC: pong %q\n", server)
					c.Write(fmt.Sprintf("PONG %s", server))
				case "INVITE":
					room := m.Trailing()
					log.Printf("IRC: joining %q\n", room)
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
					for _, p := range *plugins {
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
						log.Printf("IRC: sending: %q to %q\n", resp, to)
						c.WriteMessage(&irc.Message{
							Command: "PRIVMSG",
							Params: []string{
								to,
								resp,
							},
						})
					}
				default:
					log.Printf("IRC: unhandled - %q", m.String())
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
