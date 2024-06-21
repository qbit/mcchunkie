package chats

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"suah.dev/mcchunkie/plugins"
)

type mmail struct {
	smtpUser   string
	smtpServer string
	password   string

	imapClient *client.Client
	updateChan chan client.Update
}

func (m *mmail) buildFancyReply(msgID, to, from, originalSubject, resp string) error {
	buf := new(bytes.Buffer)
	w, err := mail.CreateWriter(buf, mail.HeaderFromMap(map[string][]string{
		"From":        {to},
		"To":          {from},
		"Subject":     {"Re: " + originalSubject},
		"References":  {msgID},
		"In-Reply-To": {msgID},
	}))
	if err != nil {
		return err
	}

	textHeader := mail.InlineHeader{}

	textPart, err := w.CreateSingleInline(textHeader)
	if err != nil {
		return err
	}

	_, err = textPart.Write([]byte(resp + "\r\n"))
	if err != nil {
		return err
	}

	if err := textPart.Close(); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return m.send(to, from, buf.Bytes())
}

func (m *mmail) send(to, from string, data []byte) error {
	conn, err := tls.Dial("tcp", m.smtpServer, &tls.Config{})
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	host := strings.TrimRight(m.smtpServer, ":")

	sc, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println(err)
		return err
	}
	defer sc.Close()

	sc.Auth(smtp.PlainAuth("", m.smtpUser, m.password, host))

	if err := sc.Mail(to); err != nil {
		log.Println(err)
		return err

	}

	if err := sc.Rcpt(from); err != nil {
		log.Println(err)
		return err
	}

	wc, err := sc.Data()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Fprint(wc, string(data))

	err = wc.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	return sc.Quit()
}

func (m *mmail) buildReply(msgID, subj, to, from, resp string) error {

	log.Printf("Mail: smtp sending mail to: %q, from: %q", from, to)

	reSubj := fmt.Sprintf("Re: %s", subj)

	wc := new(bytes.Buffer)
	fmt.Fprintf(wc, fmt.Sprintf("To: %s\r\n", from))
	fmt.Fprintf(wc, fmt.Sprintf("From: %s\r\n", to))
	fmt.Fprintf(wc, fmt.Sprintf("Subject: %s\r\n", reSubj))
	fmt.Fprintf(wc, fmt.Sprintf("References: %s\r\n", msgID))
	fmt.Fprintf(wc, fmt.Sprintf("In-Reply-To: %s\r\n", msgID))
	fmt.Fprintf(wc, "\r\n"+resp+"\r\n")

	return m.send(to, from, wc.Bytes())
}

func MailListen(store ChatStore, plugins *plugins.Plugins) error {
	var (
		smtpUser, _   = store.Get("smtp_user")
		smtpServer, _ = store.Get("smtp_server")
		imapServer, _ = store.Get("imap_server")
		imapUser, _   = store.Get("imap_user")
		mailPass, _   = store.Get("mail_password")
	)

	m := mmail{
		smtpUser:   smtpUser,
		smtpServer: smtpServer,
		password:   mailPass,
		updateChan: make(chan client.Update),
	}

	c, err := client.DialTLS(imapServer, nil)
	if err != nil {
		return err
	}
	log.Printf("Mail: connected to %q", imapServer)

	m.imapClient = c

	defer m.imapClient.Logout()

	if err = m.imapClient.Login(imapUser, mailPass); err != nil {
		return err
	}
	log.Printf("Mail: logged in as %q", imapUser)

	_, err = m.imapClient.Select("INBOX", false)
	if err != nil {
		return err
	}

	m.imapClient.Updates = m.updateChan

	go func() {
		for update := range m.updateChan {
			switch update.(type) {
			case *client.MessageUpdate:
				log.Println("Mail: received new message")

				crit := imap.NewSearchCriteria()
				crit.WithoutFlags = []string{imap.SeenFlag}

				ids, err := m.imapClient.Search(crit)
				if err != nil {
					log.Println(err)
					continue
				}

				if len(ids) == 0 {
					log.Println(len(ids))
					continue
				}

				seq := &imap.SeqSet{}
				seq.AddNum(ids...)

				messages := make(chan *imap.Message, 10)
				done := make(chan error, 1)

				// Mark the messages as read
				bodySeq := &imap.BodySectionName{
					Peek: false,
				}

				go func() {
					done <- c.Fetch(seq, []imap.FetchItem{bodySeq.FetchItem()}, messages)
				}()

				for msg := range messages {
					fullMsg := msg.GetBody(bodySeq)
					to, from, subj, msg, msgID := "", "", "", "", ""

					mr, err := mail.CreateReader(fullMsg)
					if err != nil {
						log.Println(err)
						continue
					}

					for {
						part, err := mr.NextPart()
						if err == io.EOF {
							break
						} else if err != nil {
							log.Println(err)
							continue
						}
						maybeSubj := part.Header.Get("Subject")
						if maybeSubj != "" {
							subj = maybeSubj
						}
						maybeTo := part.Header.Get("To")
						if maybeTo != "" {
							to = maybeTo
						}
						maybeFrom := part.Header.Get("From")
						if maybeFrom != "" {
							from = maybeFrom
						}
						maybeID := part.Header.Get("Message-ID")
						if maybeID != "" {
							msgID = maybeID
						}

						switch part.Header.(type) {
						case *mail.InlineHeader:
							b, _ := io.ReadAll(part.Body)
							msg = strings.TrimSpace(string(b))
						}
						log.Printf("to: %q, from: %q, subj: %q, msg: %q", to, from, subj, msg)
					}

					if to != "" && from != "" && msg != "" && subj != "" {
						for _, p := range *plugins {
							if p.Match(from, msg) {
								log.Printf("%s: responding to '%s'", p.Name(), from)
								p.SetStore(store)

								resp := p.Process(from, msg)
								log.Println(resp)

								/*
									err := m.buildReply(msgID, subj, to, from, resp)
									if err != nil {
										log.Println(err)
										continue
									}
								*/

								err := m.buildFancyReply(msgID, to, from, subj, resp)
								if err != nil {
									log.Println(err)
									continue
								}

							}
						}
					}
				}

				if err := <-done; err != nil {
					log.Fatal(err)
				}

			}
		}
	}()

	for {
		if err := m.imapClient.Noop(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}
