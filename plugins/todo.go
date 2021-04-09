package plugins

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
)

// TodoData represents a single todo
type TodoData struct {
	ID       int       `json:"id"`
	Complete bool      `json:"complete"`
	Item     string    `json:"item"`
	Date     time.Time `json:"date"`
}

// Todos are a collection of todos
type Todos []TodoData

// Save writes a set of TodoDatas to disk via a PluginStore
func (t *Todos) Save(s PluginStore) error {
	return nil
}

// String writes the todos as text
func (t *Todos) String() string {
	var s []string
	for _, ts := range *t {
		s = append(s, fmt.Sprintf("- %d %s", ts.ID, ts.Item))
	}

	if len(s) == 0 {
		return "No items!"
	}

	return strings.Join(s, "\n")
}

// GetTodos loads an existing todo set or creates a new one.
func GetTodos(key string, s PluginStore) Todos {
	tdString, _ = s.Get(key)
	return Todos{
		{
			ID:       0,
			Complete: false,
			Item:     "Milk",
			Date:     time.Now(),
		},
		{
			ID:       1,
			Complete: false,
			Item:     "Cheese",
			Date:     time.Now(),
		},
	}
}

// Todo is our plugin type
type Todo struct {
	db PluginStore
}

// SetStore is the setup function for a plugin
func (h *Todo) SetStore(s PluginStore) {
	h.db = s
}

// Descr describes this plugin
func (h *Todo) Descr() string {
	return "Simple TODO manager."
}

// Re is what our todo matches
func (h *Todo) Re() string {
	return `(?i)^todo (add|done|list|purge)\s?(.+)?$`
}

func (h *Todo) command(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$1")
}

func (h *Todo) item(msg string) string {
	re := regexp.MustCompile(h.Re())
	return re.ReplaceAllString(msg, "$2")
}

// Match checks for "home: name?" messages
func (h *Todo) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// Process does the heavy lifting
func (h *Todo) Process(from, post string) string {
	cmd := h.command(post)
	item := h.item(post)
	var tds Todos

	if cmd != "" {
		tds = GetTodos(from, h.db)
	}

	switch cmd {
	case "list":
		return tds.String()
	case "add":
		return fmt.Sprintf("added: %q", item)
	case "done":
		return fmt.Sprintf("completed: %q", item)
	case "purge":
		return "purged."
	default:
		return "unknown command"
	}
}

// RespondText to looking up of weather lookup requests
func (h *Todo) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendMD(c, ev.RoomID, h.Process(ev.Sender, post))
}

// Name Todo!
func (h *Todo) Name() string {
	return "Todo"
}
