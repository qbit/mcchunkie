package chats

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
	"suah.dev/mcchunkie/plugins"
)

// Matrix lets us connect to Matrix chat servers
type Matrix struct {
	Client *gomatrix.Client
	Store  plugins.PluginStore
	Name   string
}

// Connect connects to a matrix server
func (m *Matrix) Connect(store plugins.PluginStore) error {
	server, err := store.Get("matrix_server")
	if err != nil {
		return err
	}

	fmt.Println(server)

	m.Client, err = gomatrix.NewClient(
		server,
		"",
		"",
	)
	if err != nil {
		return err
	}

	m.Name, _ = store.Get("matrix_user")
	accessToken, _ := store.Get("matrix_access_token")
	userID, _ := store.Get("matrix_user_id")
	botOwner, _ := store.Get("bot_owner")

	m.Client.SetCredentials(userID, accessToken)
	m.Client.Store = NewStore()
	syncer := gomatrix.NewDefaultSyncer(m.Name, m.Client.Store)
	m.Client.Client = http.DefaultClient
	m.Client.Syncer = syncer

	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if ev.Sender == m.Name {
			return
		}
		switch ev.Sender {
		case botOwner:
			if ev.Content["membership"] == "invite" {
				log.Printf("Joining %s (invite from %s)\n", ev.RoomID, ev.Sender)
				if _, err := m.Client.JoinRoom(ev.RoomID, "", nil); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	})

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == m.Name {
			return
		}

		// Sending a response per plugin hits issues, so save them and
		// send as one message.
		var helps []string
		for _, p := range plugins.Plugs {
			var post string
			var ok bool
			if post, ok = ev.Body(); !ok {
				// Invaild body, for some reason
				return
			}
			if mtype, ok := ev.MessageType(); ok {
				switch mtype {
				case "m.text":
					if p.Match(m.Name, post) {
						log.Printf("%s: responding to '%s'", p.Name(), ev.Sender)
						p.SetStore(store)

						err := m.Process(&Message{
							From: m.Name,
							Body: post,
						})
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
		if len(helps) > 0 {
			err := plugins.SendMD(m.Client, ev.RoomID, strings.Join(helps, "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	})

	go func() {
		log.Println("MATRIX: syncing..")
		if err := m.Client.Sync(); err != nil {
			fmt.Println("Sync() returned ", err)
		}

		time.Sleep(1 * time.Second)
	}()

	return nil
}

// Disconnect disconnects cleanely from our service
func (m *Matrix) Disconnect() error {
	m.Client.StopSync()
	return nil
}

// Process receives a Message and determines if it should be used or not
func (m *Matrix) Process(msg *Message) error {
	if msg.Service != "Matrix" {
		return nil
	}
	for _, p := range plugins.Plugs {
		if p.Match(m.Name, msg.Body) {
			log.Printf("%s: responding to '%s'", p.Name(), msg.From)
			p.SetStore(m.Store)
			resp := p.Process(msg.From, msg.Body)
			log.Println(resp)
		}
	}
	return nil
}
