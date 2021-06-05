package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"suah.dev/mcchunkie/chats"
	"suah.dev/mcchunkie/plugins"
	"suah.dev/protect"
)

const header = `
# mcchunkie

[![mcchunkie's face](mcchunkie.png)](https://builds.sr.ht/~qbit/mcchunkie?)

[![builds.sr.ht status](https://builds.sr.ht/~qbit/mcchunkie.svg)](https://builds.sr.ht/~qbit/mcchunkie?)

A [Matrix](https://matrix.org) chat bot.`

func main() {
	var username, db, avatar, prof string
	var key, value, get string
	var setup, doc, verbose bool

	flag.BoolVar(&doc, "doc", false, "print plugin information and exit")
	flag.BoolVar(&setup, "s", false, "setup account")
	flag.BoolVar(&verbose, "v", false, "print verbose messages")

	flag.StringVar(&avatar, "avatar", "", "set the avatar of the bot to specified url")
	flag.StringVar(&db, "db", "db", "full path to database directory")
	flag.StringVar(&get, "get", "", "grab an entry from the store")
	flag.StringVar(&key, "key", "", "create an entry in the data store listed under 'key'")
	flag.StringVar(&username, "user", "", "username to connect to matrix server with")
	flag.StringVar(&value, "value", "", "set the value of 'key' to be stored")
	flag.StringVar(&prof, "prof", "", "listen string for pprof")

	flag.Parse()

	_ = protect.Pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	_ = protect.Unveil("/etc/resolv.conf", "r")
	_ = protect.Unveil("/etc/ssl/cert.pem", "r")
	_ = protect.Unveil(db, "rwc")

	var err = protect.UnveilBlock()
	if err != nil {
		log.Fatal(err)
	}

	store, err := NewStore(db)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if key != "" && value != "" {
		store.Set(key, value)
		os.Exit(0)
	}

	if doc {
		fmt.Println(header)
		fmt.Println("\n|Plugin Name|Match|Description|")
		fmt.Println("|----|---|---|")
		for _, p := range plugins.Plugs {
			fmt.Printf("|%s|`%s`|%s|\n", p.Name(), strings.ReplaceAll(p.Re(), "|", "\\|"), p.Descr())
		}
		os.Exit(0)
	}

	if get != "" {
		val, err := store.Get(get)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		fmt.Println(val)
		os.Exit(0)
	}

	go func() {
		for _, c := range chats.ChatMethods {
			c.Connect(store)
		}
	}()

	defer func() {
		for _, c := range chats.ChatMethods {
			c.Disconnect()
		}
	}()

	for {
		time.Sleep(time.Second * 60)
	}
}
