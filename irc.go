package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"gopkg.in/irc.v3"
	"suah.dev/mcchunkie/plugins"
)

func ircConnect(store *FStore, plugins *plugins.Plugins) error {
	var ircServer, _ = store.Get("irc_server")
	var ircPort, _ = store.Get("irc_port")
	var ircNick, _ = store.Get("irc_nick")
	var ircPass, _ = store.Get("irc_pass")
	var ircRooms, _ = store.Get("irc_rooms")
	//var toRE = regexp.MustCompile(`^:(\w+)\s`)

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
			log.Fatalln(err)
		}
	}
	return nil
}
