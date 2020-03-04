package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
	"suah.dev/mcchunkie/plugins"
)

func main() {
	var username, password, userID, accessToken, server, db, avatar, botOwner string
	var key, value, get string
	var setup, doc, verbose bool

	flag.BoolVar(&doc, "doc", false, "print plugin information and exit")
	flag.BoolVar(&setup, "s", false, "setup account")
	flag.BoolVar(&verbose, "v", false, "print verbose messages")

	flag.StringVar(&avatar, "avatar", "", "set the avatar of the bot to specified url")
	flag.StringVar(&db, "db", "db", "full path to database directory")
	flag.StringVar(&get, "get", "", "grab an entry from the store")
	flag.StringVar(&key, "key", "", "create an entry in the data store listed under 'key'")
	flag.StringVar(&server, "server", "", "matrix server")
	flag.StringVar(&username, "user", "", "username to connect to matrix server with")
	flag.StringVar(&value, "value", "", "set the value of 'key' to be stored")

	flag.Parse()

	pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	unveil("/etc/resolv.conf", "r")
	unveil("/etc/ssl/cert.pem", "r")
	unveil(db, "rwc")
	unveilBlock()

	var help = `^help: (\w+)$`
	var helpRE = regexp.MustCompile(help)
	var store, err = NewStore(db)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if key != "" && value != "" {
		store.Set(key, value)
		os.Exit(0)
	}

	if doc {
		fmt.Println("|Plugin Name|Match|Description|")
		fmt.Println("|----|---|---|")
		for _, p := range plugins.Plugs {
			fmt.Printf("|%s|`%s`|%s|\n", p.Name(), p.Re(), p.Descr())
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

	if server == "" {
		server, err = store.Get("server")
		if err != nil {
			if err != nil {
				log.Fatalf("%s\n", err)
			}
		}
		if server == "" {
			log.Fatalln("please specify a server")
		}

	} else {
		store.Set("server", server)
	}

	log.Printf("connecting to %s\n", server)

	cli, err := gomatrix.NewClient(
		server,
		"",
		"",
	)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if setup {
		log.Println("requesting access token")
		password, err = prompt(fmt.Sprintf("Password for '%s': ", username))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println()

		resp, err := cli.Login(&gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: password,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// No longer need tty now that we have our info
		pledge("stdio unveil rpath wpath cpath flock dns inet")

		store.Set("username", username)
		store.Set("access_token", resp.AccessToken)
		store.Set("user_id", resp.UserID)

		accessToken = resp.AccessToken
		userID = resp.UserID
	} else {
		username, _ = store.Get("username")
		accessToken, _ = store.Get("access_token")
		userID, _ = store.Get("user_id")
		botOwner, _ = store.Get("bot_owner")
	}

	cli.SetCredentials(userID, accessToken)
	cli.Store = store
	syncer := gomatrix.NewDefaultSyncer(username, store)
	cli.Client = http.DefaultClient
	cli.Syncer = syncer

	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}

		if ev.Sender == botOwner && ev.Content["membership"] == "invite" {
			log.Printf("Joining %s (invite from %s)\n", ev.RoomID, ev.Sender)
			if _, err := cli.JoinRoom(ev.RoomID, "", nil); err != nil {
				log.Fatalln(err)
			}
		}
	})

	go func() {
		errataCount := 0
		storeCount, _ := store.Get("errata_count")
		openbsdRelease, _ := store.Get("openbsd_release")
		errataCount, err = strconv.Atoi(storeCount)
		for {
			got, err := ParseRemoteErrata(
				fmt.Sprintf("https://www.openbsd.org/errata%s.html", openbsdRelease),
			)
			if err != nil {
				fmt.Println(err)
			}
			l := got.Length
			if l > errataCount {
				alertRooms, _ := store.Get("errata_rooms")
				c := 0
				for _, errata := range got.List {
					if c+1 > errataCount {
						log.Printf("%03d: %s - %s\n", errata.ID, errata.Type, errata.Desc)
						for _, room := range strings.Split(alertRooms, ",") {
							cli.SendNotice(room, PrintErrata(&errata))
						}
					}
					c = c + 1
				}
				errataCount = l
			}
			store.Set("errata_count", strconv.Itoa(l))
			time.Sleep(2 * time.Hour)
		}
	}()

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}

		// Sending a response per plugin hits issues, so save them and
		// send as one message.
		helps := []string{}
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
					if helpRE.Match([]byte(post)) {
						pn := p.Name()
						hName := helpRE.ReplaceAllString(post, "$1")
						if hName == pn || hName == "puke" {
							helps = append(helps, fmt.Sprintf("**%s**: `%s` -  _%s_\n", p.Name(), p.Re(), p.Descr()))
						}

					}
					if p.Match(username, post) {
						log.Printf("%s: responding to '%s'", p.Name(), ev.Sender)
						p.SetStore(store)

						start := time.Now()
						p.RespondText(cli, ev, username, post)
						elapsed := time.Since(start)
						if verbose {
							log.Printf("%s took %s to run\n", p.Name(), elapsed)
						}
					}
				}
			}
		}
		if len(helps) > 0 {
			plugins.SendMD(cli, ev.RoomID, strings.Join(helps, "\n"))
		}
	})

	if avatar != "" {
		log.Printf("Setting avatar to: '%s'", avatar)
		rmu, err := cli.UploadLink(avatar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = cli.SetAvatarURL(rmu.ContentURI)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	for {
		log.Println("syncing..")
		if err := cli.Sync(); err != nil {
			fmt.Println("Sync() returned ", err)
		}

		time.Sleep(1 * time.Second)
	}
}
