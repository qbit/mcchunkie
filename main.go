package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"git.sr.ht/~qbit/mcchunkie/plugins"
	"github.com/matrix-org/gomatrix"
)

func messageToMe(sn, message string) bool {
	return strings.Contains(message, sn)
}

func sendMessage(c *gomatrix.Client, roomID, message string) error {
	_, err := c.UserTyping(roomID, true, 3)
	if err != nil {
		return err
	}

	c.SendText(roomID, message)

	_, err = c.UserTyping(roomID, false, 0)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var username, password, userID, accessToken, server, db, avatar string
	var setup bool

	flag.StringVar(&username, "user", "", "username to connect to matrix server with")
	flag.StringVar(&server, "server", "", "matrix server")
	flag.StringVar(&avatar, "avatar", "", "set the avatar of the bot to specified url")
	flag.BoolVar(&setup, "s", false, "setup account")
	flag.StringVar(&db, "db", "db", "full path to database directory")

	flag.Parse()

	pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	unveil("/etc/resolv.conf", "r")
	unveil("/etc/ssl/cert.pem", "r")
	unveil(db, "rwc")
	unveilBlock()

	var store, err = NewStore(db)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if server == "" {
		server, err = store.get("server")
		if server == "" {
			log.Fatalln("please specify a server")
		}

	} else {
		store.set("server", server)
	}

	log.Printf("connecting to %s\n", server)

	cli, err := gomatrix.NewClient(
		server,
		"",
		"",
	)

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

		store.set("username", username)
		store.set("access_token", resp.AccessToken)
		store.set("user_id", resp.UserID)

		accessToken = resp.AccessToken
		userID = resp.UserID
	} else {
		username, _ = store.get("username")
		accessToken, _ = store.get("access_token")
		userID, _ = store.get("user_id")
	}

	cli.SetCredentials(userID, accessToken)
	cli.Store = store
	syncer := gomatrix.NewDefaultSyncer(username, store)
	cli.Client = http.DefaultClient
	cli.Syncer = syncer

	/*
		// TODO: Add ability to join / part rooms

		if _, err := cli.JoinRoom("!tmCVBJAeuKjCfihUjb:cobryce.com", "", nil); err != nil {
			log.Fatalln(err)
		}
		if _, err := cli.JoinRoom("!sFPUeGfHqjiItcjNIN:matrix.org", "", nil); err != nil {
			log.Fatalln(err)
		}
		if _, err := cli.JoinRoom("!ALCZnrYadLGSySIFZr:matrix.org", "", nil); err != nil {
			log.Fatalln(err)
		}
		if _, err := cli.JoinRoom("!LTxJpLHtShMVmlpwmZ:tapenet.org", "", nil); err != nil {
			log.Fatalln(err)
		}
		if _, err := cli.JoinRoom("!TjjamgVanKpNiswkoJ:pintobyte.com", "", nil); err != nil {
			log.Fatalln(err)
		}
	*/

	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if ev.Sender == username {
			return
		}

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
					p.RespondText(cli, ev, username, post)
				}
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
