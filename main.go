package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"suah.dev/mcchunkie/chats"
	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
	"suah.dev/protect"
)

const header = `
# mcchunkie

[![mcchunkie's face](mcchunkie.png)](https://builds.sr.ht/~qbit/mcchunkie?)

[![builds.sr.ht status](https://builds.sr.ht/~qbit/mcchunkie.svg)](https://builds.sr.ht/~qbit/mcchunkie?)

A Matrix, XMPP, IRC, Mail and SMS chat bot.`

func main() {
	var db string
	var key, value, get, disableChats, disablePlugins string
	var doc bool

	flag.BoolVar(&doc, "doc", false, "print plugin information and exit")
	flag.StringVar(&db, "db", "db", "full path to database directory")
	flag.StringVar(&get, "get", "", "grab an entry from the store")
	flag.StringVar(&key, "key", "", "create an entry in the data store listed under 'key'")
	flag.StringVar(&value, "value", "", "set the value of 'key' to be stored")
	flag.StringVar(&disableChats, "dc", "", fmt.Sprintf("comma delimited list of chat types to disable (case insensitive)\nEnabled by default: %s", chats.ChatMethods.List()))
	flag.StringVar(&disablePlugins, "dp", "", fmt.Sprintf("comma delimited list of plugin types to disable (case insensitive)\nEnabled by default: %s", plugins.Plugs.List()))

	flag.Parse()

	_ = protect.Pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	_ = protect.Unveil("/etc/resolv.conf", "r")
	_ = protect.Unveil("/etc/ssl/cert.pem", "r")
	_ = protect.Unveil(db, "rwc")

	var err = protect.UnveilBlock()
	if err != nil {
		log.Fatal(err)
	}

	store, err := mcstore.NewStore(db)
	if err != nil {
		log.Fatalln(err)
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

	disableList := strings.Split(strings.ToLower(disableChats), ",")
	chatEnabled := func(chat string) bool {
		if slices.Contains(disableList, strings.ToLower(chat)) {
			return false
		}
		return true
	}
	for _, chat := range chats.ChatMethods {
		go func() {
			if chatEnabled(chat.Name()) {
				for {
					log.Printf("Starting %s...", chat.Name())
					err := chat.Connect(store, &plugins.Plugs)
					if err != nil {
						log.Println(fmt.Errorf("%s: %q", chat.Name(), err))
					}
					time.Sleep(15 * time.Second)
				}
			}
		}()
	}

	gotChat, err := chats.ChatMethods.ByName("IRC")
	go chats.GotListen(store, gotChat)

	for {
		errataCount := 0
		storeCount, err := store.Get("errata_count")
		if err != nil {
			log.Fatal(err)
		}
		openbsdRelease, err := store.Get("openbsd_release")
		if err != nil {
			log.Fatal(err)
		}
		errataCount, err = strconv.Atoi(storeCount)

		got, err := ParseRemoteErrata(
			fmt.Sprintf("http://ftp.openbsd.org/pub/OpenBSD/patches/%s/common/",
				openbsdRelease,
			),
		)
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Hour)
			continue
		}
		l := len(got.List)
		if l > errataCount {
			alertRooms, err := store.Get("errata_rooms")
			if err != nil {
				log.Fatal(err)
			}
			c := 0
			for _, erratum := range got.List {
				if c+1 > errataCount {
					log.Printf("Notifying for erratum %03d\n", erratum.ID)
					err = erratum.Fetch()
					if err != nil {
						fmt.Println(err)
						break
					}
					for _, room := range strings.Split(alertRooms, ",") {
						for _, c := range chats.ChatMethods {
							if chatEnabled(c.Name()) {
								err = c.Send(room, PrintErrataMD(&erratum))
								if err != nil {
									log.Printf("%s: %q", c.Name(), err)
								}
							}
						}
					}
				}
				c = c + 1
			}
			errataCount = l
		}
		store.Set("errata_count", strconv.Itoa(l))
		time.Sleep(2 * time.Hour)
	}
}
