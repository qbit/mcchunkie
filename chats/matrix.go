package chats

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matrix-org/gomatrix"
	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
)

type MatrixChat struct {
	client *gomatrix.Client
}

func (mc *MatrixChat) Name() string { return "Matrix" }

func (mc *MatrixChat) Send(to, msg string) error {
	return plugins.SendMDNotice(mc.client, to, msg)
}

func (mc *MatrixChat) Connect(store *mcstore.MCStore, plugs *plugins.Plugins) error {
	server, err := store.Get("matrix_server")
	if err != nil {
		return err
	}
	log.Printf("Matrix: connecting to %s\n", server)

	mc.client, err = gomatrix.NewClient(
		server,
		"",
		"",
	)
	if err != nil {
		return err
	}

	username, err := store.Get("matrix_username")
	if err != nil {
		return err
	}
	accessToken, err := store.Get("matrix_access_token")
	if err != nil {
		return err
	}
	userID, err := store.Get("matrix_user_id")
	if err != nil {
		return err
	}
	botOwner, err := store.Get("matrix_bot_owner")
	if err != nil {
		return err
	}

	mc.client.SetCredentials(userID, accessToken)
	mc.client.Store = store
	syncer := gomatrix.NewDefaultSyncer(username, store)
	mc.client.Client = http.DefaultClient
	mc.client.Syncer = syncer

	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}
		switch ev.Sender {
		case botOwner:
			if ev.Content["membership"] == "invite" {
				log.Printf("Joining %s (invite from %s)\n", ev.RoomID, ev.Sender)
				if _, err := mc.client.JoinRoom(ev.RoomID, "", nil); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	})

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}

		for _, p := range *plugs {
			var post string
			var ok bool
			if post, ok = ev.Body(); !ok {
				return
			}
			if mtype, ok := ev.MessageType(); ok {
				switch mtype {
				case "m.text":
					if p.Match(username, post) {
						log.Printf("%s: responding to '%s'", p.Name(), ev.Sender)
						p.SetStore(store)

						err := p.RespondText(mc.client, ev, username, post)
						if err != nil {
							fmt.Println(err)
							plugins.SendText(mc.client, ev.RoomID, err.Error())
						}
					}
				}
			}
		}
	})

	return mc.client.Sync()
}
