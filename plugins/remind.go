package plugins

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/matrix-org/gomatrix"
)

// Remind responds to remind requests
type Remind struct {
}

type Reminder struct {
	Duration time.Duration
	String   string
}

// Descr describes this plugin
func (h *Remind) Descr() string {
	return "Remind lets one set reminders to be ping'd on at a later point in time"
}

// Re returns the remind matching string
func (h *Remind) Re() string {
	return `(?i)^remind: (?P<duration>(\w+)) (?P<reminder>(.+))$`
}

func (h *Remind) fix(msg string) (*Reminder, error) {
	re := regexp.MustCompile(h.Re())
	r := &Reminder{}

	matches := re.FindStringSubmatch(msg)
	for i, name := range re.SubexpNames() {
		var err error
		switch name {
		case "duration":
			log.Println("duration", matches[i])
			r.Duration, err = time.ParseDuration(matches[i])
			if err != nil {
				return nil, err
			}
		case "reminder":
			log.Println("reminder", matches[i])
			r.String = matches[i]
		}
	}

	log.Println(r)

	re.ReplaceAllString(msg, "$1")
	return r, nil
}

// Match determines if we should call the response for Remind
func (h *Remind) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here.
func (h *Remind) SetStore(_ PluginStore) {}

func (h *Remind) Process(from, msg string) (string, func() string) {
	r, err := h.fix(msg)
	if err != nil {
		return err.Error(), func() string { return "" }
	}
	now := time.Now()
	resp := fmt.Sprintf("OK %s, I'll remind you on %s", from, now.Add(r.Duration).Format(time.RFC1123))

	return resp, func() string {
		r := r
		ch := make(chan string)
		go func() {
			time.Sleep(r.Duration)
			ch <- fmt.Sprintf("%s: %s", from, r.String)
		}()
		return <-ch
	}
}

// RespondText to looking up of remind requests
func (h *Remind) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	resp, delayedResp := h.Process(ev.Sender, post)
	go func() {
		SendText(c, ev.RoomID, delayedResp())
	}()

	return SendText(c, ev.RoomID, resp)
}

// Name Remind!
func (h *Remind) Name() string {
	return "Remind"
}
