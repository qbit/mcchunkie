package plugins

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/matrix-org/gomatrix"
	"github.com/ollama/ollama/api"
)

// Llama responds to llama messages
type Llama struct {
	client *api.Client
	db     PluginStore
}

func (l *Llama) Descr() string {
	return "Send queries to a local instance of ollama"
}

func (l *Llama) Re() string {
	return `(?i)^o?llama:(.+)$`
}

func (l *Llama) Match(_, msg string) bool {
	re := regexp.MustCompile(l.Re())
	return re.MatchString(msg)
}

func (l *Llama) SetStore(s PluginStore) {
	l.db = s
}

func (l *Llama) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, delayedResp := l.Process(ev.Sender, post)
	go func() {
		SendText(c, ev.RoomID, delayedResp())
	}()

	return SendMD(c, ev.RoomID, resp)
}

func (l *Llama) Process(from, msg string) (string, func() string) {
	var err error
	ctx := context.Background()

	re := regexp.MustCompile(l.Re())
	query := re.ReplaceAllString(msg, "$1")
	llamaServer, err := l.db.Get("ollama_host")
	if err != nil {
		return err.Error(), RespStub
	}

	botOwners, err := l.db.Get("bot_owners")
	if err != nil {
		return err.Error(), RespStub
	}
	log.Println("owners, from", botOwners, from)
	if !slices.Contains(strings.Split(botOwners, ","), from) {
		return errors.New(fmt.Sprintf("sorry, %s, I can't let you do that.", from)).Error(), RespStub
	}

	if l.client == nil {
		u, err := url.Parse(llamaServer)
		if err != nil {
			return err.Error(), RespStub
		}
		l.client = api.NewClient(u, http.DefaultClient)
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: "provide very brief, concise, single line responses unless asked otherwise",
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := &api.ChatRequest{
		Model:    "llama3.1",
		Messages: messages,
	}

	respSet := []string{}
	err = l.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		respSet = append(respSet, resp.Message.Content)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return strings.Join(respSet, ""), RespStub
}

func (l *Llama) Name() string {
	return "Llama"
}
