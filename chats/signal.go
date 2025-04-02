package chats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"regexp"

	"suah.dev/mcchunkie/mcstore"
	"suah.dev/mcchunkie/plugins"
)

type SigMsg struct {
	To  string
	Msg string
}

type GroupInfo struct {
	GroupID string `json:"groupId"`
}

type Mention struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type DataMessage struct {
	Recipient []string  `json:"recipient,omitempty"`
	ID        string    `json:"id"`
	GroupID   string    `json:"groupId,omitempty"`
	Timestamp int64     `json:"timestamp"`
	Message   string    `json:"message"`
	GroupInfo GroupInfo `json:"groupInfo"`
	Mentions  []Mention `json:"mentions,omitempty"`
}

type Envelope struct {
	SourceNumber string      `json:"sourceNumber"`
	SourceUUID   string      `json:"sourceUuid"`
	DataMessage  DataMessage `json:"dataMessage"`
}

type Event struct {
	Envelope Envelope `json:"envelope"`
}

type SignalChat struct {
	number string
	socket string
	in     chan []byte
	out    chan []byte
}

type SendEvent struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  DataMessage `json:"params"`
}

type ReceiveEvent struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Event  `json:"params"`
}

func randID() string {
	char := "abcdefghijklmnopqrstuvwxyz1234567890"
	b := make([]byte, 16)
	for i := range b {
		b[i] = char[rand.Intn(len(char))]
	}
	return string(b)
}

func NewSendEvent() *SendEvent {
	se := &SendEvent{
		JSONRPC: "2.0",
		Method:  "send",
	}
	return se
}

func (x *SignalChat) Send(to string, resp string) error {
	se := NewSendEvent()
	se.Params.Message = resp
	se.Params.ID = randID()

	uuidRE := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if uuidRE.Match([]byte(to)) {
		se.Params.Recipient = append(se.Params.Recipient, to)
	} else {
		se.Params.GroupID = to
	}

	data, err := json.Marshal(se)
	if err != nil {
		return err
	}

	x.in <- data

	return nil
}

func (x *SignalChat) Name() string {
	return "Signal"
}

func (x *SignalChat) Connect(store *mcstore.MCStore, plugins *plugins.Plugins) error {
	number, _ := store.Get("signal_number")
	socket, _ := store.Get("signal_socket")
	if x.number == "" {
		x.number = number
		x.socket = socket
	}

	c, err := net.Dial("unix", x.socket)
	if err != nil {
		log.Println("Signal", err)
		return err
	}

	defer c.Close()

	x.in = make(chan []byte)
	x.out = make(chan []byte)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := c.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Signal", err)
				}
				close(x.in)
				return
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			x.out <- data
		}
	}()

	go func() {
		for data := range x.in {
			_, err := c.Write(append(data, []byte("\n")...))
			if err != nil {
				log.Println("Signal", err)
				return
			}
		}
	}()

	for {
		select {
		case data, ok := <-x.out:
			if !ok {
				return fmt.Errorf("disconnected")
			}
			events := []ReceiveEvent{}
			for _, ev := range bytes.Split(data, []byte("\n")) {
				e := ReceiveEvent{}
				if len(ev) == 0 {
					continue
				}
				log.Println("RAW", string(ev))
				err = json.Unmarshal(ev, &e)
				if err != nil {

					continue
				}
				if e.Method == "receive" {
					events = append(events, e)
				}

				for _, event := range events {
					if event.Params.Envelope.DataMessage.Message == "" {
						continue
					}

					from := event.Params.Envelope.SourceUUID
					msg := event.Params.Envelope.DataMessage.Message

					if event.Params.Envelope.DataMessage.GroupInfo.GroupID != "" {
						from = event.Params.Envelope.DataMessage.GroupInfo.GroupID
					}

					resp := ""
					delayedResp := func() string { return "" }

					for _, p := range *plugins {

						if p.Match(from, msg) {
							p.SetStore(store)
							resp, delayedResp = p.Process(from, msg)
						}
					}
					if resp != "" {
						log.Printf("Signal: sending: %q to %q\n", resp, from)
						x.Send(from, resp)
						go func() {
							dresp := delayedResp()
							if dresp != "" {
								log.Printf("Signal: sending: %q to %q\n", dresp, from)
								x.Send(from, dresp)
							}
						}()
					}
				}
			}
		}
	}
}
