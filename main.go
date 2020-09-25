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
	"golang.org/x/crypto/bcrypt"
	"suah.dev/mcchunkie/plugins"
	"suah.dev/protect"
)

const header = `
# mcchunkie

[![builds.sr.ht status](https://builds.sr.ht/~qbit/mcchunkie.svg)](https://builds.sr.ht/~qbit/mcchunkie?)`

func main() {
	var username, shortName, password, userID, accessToken, server, db, avatar, botOwner string
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

	_ = protect.Pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	_ = protect.Unveil("/etc/resolv.conf", "r")
	_ = protect.Unveil("/etc/ssl/cert.pem", "r")
	_ = protect.Unveil(db, "rwc")

	var err = protect.UnveilBlock()
	if err != nil {
		log.Fatal(err)
	}

	var help = `^help: (\w+)$`
	var helpRE = regexp.MustCompile(help)
	var kvRE = regexp.MustCompile(`^(.+)\s->\s(.+)$`)
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
		_ = protect.Pledge("stdio unveil rpath wpath cpath flock dns inet")

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

	shortName = plugins.NameRE.ReplaceAllString(username, "$1")

	cli.SetCredentials(userID, accessToken)
	cli.Store = store
	syncer := gomatrix.NewDefaultSyncer(username, store)
	cli.Client = http.DefaultClient
	cli.Syncer = syncer

	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}
		switch ev.Sender {
		case botOwner:
			if ev.Content["membership"] == "invite" {
				log.Printf("Joining %s (invite from %s)\n", ev.RoomID, ev.Sender)
				if _, err := cli.JoinRoom(ev.RoomID, "", nil); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	})

	go func() {
		var htpass, _ = store.Get("got_htpass")
		var got_room, _ = store.Get("got_room")
		var got_port, _ = store.Get("got_listen")

		log.Printf("GOT: listening on %q and sending messages to %q\n", got_port, got_room)

		http.HandleFunc("/_got", func(w http.ResponseWriter, r *http.Request) {
			var msg string
			user, pass, ok := r.BasicAuth()
			err := bcrypt.CompareHashAndPassword([]byte(htpass), []byte(pass))
			if !(ok && err == nil && user == "got") {
				log.Printf("GOT: failed auth '%s'\n", user)
				w.Header().Set("WWW-Authenticate", `Basic realm="got notify"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = r.ParseForm()
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}

			switch r.Method {
			case http.MethodGet:
				msg = r.Form.Get("message")
			case http.MethodPost:
				msg = r.Form.Get("file")
			default:
				http.Error(w, fmt.Sprintf("method %q not implemented", r.Method), http.StatusMethodNotAllowed)
				return
			}

			log.Printf("GOT: sending '%s'\n", msg)
			err = plugins.SendMDNotice(cli, got_room, msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("can not send commit info: %s", err), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "ok")

		})

		log.Fatal(http.ListenAndServe(got_port, nil))
	}()

	go func() {
		for {
			errataCount := 0
			storeCount, _ := store.Get("errata_count")
			openbsdRelease, _ := store.Get("openbsd_release")
			errataCount, err = strconv.Atoi(storeCount)

			got, err := ParseRemoteErrata(
				fmt.Sprintf("http://ftp.openbsd.org/pub/OpenBSD/patches/%s/common/", openbsdRelease),
			)
			if err != nil {
				fmt.Println(err)
				time.Sleep(2 * time.Hour)
				continue
			}
			l := len(got.List)
			if l > errataCount {
				alertRooms, _ := store.Get("errata_rooms")
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
							err = plugins.SendMDNotice(cli, room, PrintErrataMD(&erratum))
							if err != nil {
								fmt.Println(err)
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
	}()

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}

		switch ev.Sender {
		case botOwner:
			var post string
			var ok bool

			if post, ok = ev.Body(); !ok {
				return
			}

			if plugins.ToMe(username, post) {
				mp := plugins.RemoveName(shortName, post)
				if kvRE.MatchString(mp) {
					key := kvRE.ReplaceAllString(mp, "$1")
					val := kvRE.ReplaceAllString(mp, "$2")
					store.Set(key, val)
					log.Printf("Setting %q to %q", key, val)
					err := plugins.SendMD(cli, ev.RoomID, fmt.Sprintf("Set **%q** = *%q*", key, val))
					if err != nil {
						log.Println(err)
					}
					return
				}
			}
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
						err := p.RespondText(cli, ev, username, post)
						if err != nil {
							fmt.Println(err)
						}
						elapsed := time.Since(start)
						if verbose {
							log.Printf("%s took %s to run\n", p.Name(), elapsed)
						}
					}
				}
			}
		}
		if len(helps) > 0 {
			err := plugins.SendMD(cli, ev.RoomID, strings.Join(helps, "\n"))
			if err != nil {
				log.Println(err)
			}
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
