package chats

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/matrix-org/gomatrix"
	"suah.dev/mcchunkie/plugins"
)

// Matrix lets us connect to Matrix chat servers
type Matrix struct {
	Client *gomatrix.Client
	Store  plugins.PluginStore
	Name   string
	Log    *log.Logger
}

// Connect connects to a matrix server
func (m *Matrix) Connect(store plugins.PluginStore) error {
	server, err := store.Get("matrix_server")
	if err != nil {
		return err
	}

	m.Log = log.Default()
	m.Log.SetPrefix("MATRIX: ")

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
				m.Log.Printf("Joining %s (invite from %s)\n", ev.RoomID, ev.Sender)
				if _, err := m.Client.JoinRoom(ev.RoomID, "", nil); err != nil {
					m.Log.Fatalln(err)
				}
				return
			}
		}
	})

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == m.Name {
			return
		}

		var post string
		var ok bool
		if post, ok = ev.Body(); !ok {
			// Invaild body, for some reason
			return
		}
		if mtype, ok := ev.MessageType(); ok {
			switch mtype {
			case "m.text":
				err := m.Process(&Message{
					Service: "Matrix",
					Sender:  ev.Sender,
					Body:    post,
					Room:    ev.RoomID,
				})
				if err != nil {
					m.Log.Println(err)
				}
			}
		}
	})

	go func() {
		m.Log.Println("syncing..")
		if err := m.Client.Sync(); err != nil {

			m.Log.Printf("Sync() returned: \n")
			m.Log.Println(err)
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
	m.Log.Printf("%#v\n", msg)
	if msg.Service != "Matrix" {
		return nil
	}

	for _, p := range plugins.Plugs {
		if p.Match(m.Name, msg.Body) {
			m.Log.Printf("%s: responding to '%s'", p.Name(), msg.Sender)
			p.SetStore(m.Store)
			_, err := m.Client.UserTyping(msg.Room, true, 3)
			if err != nil {
				return err
			}

			_, err = m.Client.SendText(msg.Room, p.Process(msg.Sender, msg.Body))
			if err != nil {
				return err
			}

			_, err = m.Client.UserTyping(msg.Room, false, 0)
			if err != nil {
				return err
			}

			return nil
		}
	}
	return nil
}
