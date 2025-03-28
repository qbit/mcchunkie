package chats

import (
	// "github.com/agl/xmpp-client/xmpp"

	"log"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
)

type XMPPChat struct {
}

func (x *XMPPChat) Send(string, string) error {
	return nil
}

func (x *XMPPChat) Name() string {
	return "XMPP"
}

// XMPPConnect connects to our irc server
func (x *XMPPChat) Connect(store *mcstore.MCStore, plugins *plugins.Plugins) error {
	jid, _ := store.Get("xmpp_jid")
	pass, _ := store.Get("xmpp_pass")
	server, _ := store.Get("xmpp_server")

	config := &xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: server,
		},
		Jid:        jid,
		Credential: xmpp.Password(pass),
	}
	log.Printf("XMPP: connecting to %q", server)

	router := xmpp.NewRouter()
	router.HandleFunc("message", func(s xmpp.Sender, p stanza.Packet) {
		msg, ok := p.(stanza.Message)
		if !ok {
			return
		}

		resp := ""
		delayedResp := func() string { return "" }
		log.Println(msg.From, msg.To)
		log.Println(msg.Body)
		for _, p := range *plugins {
			if p.Match(msg.From, msg.Body) {
				log.Printf("%s: responding to '%s'", p.Name(), msg.From)
				p.SetStore(store)
				resp, delayedResp = p.Process(msg.From, msg.Body)
			}
		}
		if resp != "" {
			log.Printf("XMPP: sending: %q to %q\n", resp, msg.From)
			reply := stanza.Message{Attrs: stanza.Attrs{To: msg.From}, Body: resp}
			_ = s.Send(reply)
			go func() {
				dresp := delayedResp()
				if dresp != "" {
					log.Printf("XMPP: sending: %q to %q\n", dresp, msg.From)
					reply := stanza.Message{Attrs: stanza.Attrs{To: msg.From}, Body: dresp}
					_ = s.Send(reply)
				}
			}()

		}
	})

	client, err := xmpp.NewClient(config, router, func(err error) {
		log.Printf("XMPP: %q", err)
	})
	if err != nil {
		return err
	}

	cm := xmpp.NewStreamManager(client, nil)
	return cm.Run()
}
