package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/chats"
	"suah.dev/mcchunkie/mcstore"
)

func main() {
	td, err := os.MkdirTemp("", "got-test")
	if err != nil {
		log.Fatal(err)
	}

	store, err := mcstore.NewStore(td)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("testing1234"), 11)
	if err != nil {
		log.Fatal(err)
	}

	store.Set("got_listen", ":8043")
	store.Set("got_htpass", string(hash))
	store.Set("got_room", "stdout")

	chat := &chats.IRCChat{}

	chats.GotListen(store, chat)
}
