package plugins

import (
	"fmt"
	"image"
	"image/color"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// Palette responds to color messages
type Palette struct {
}

// Descr describes this plugin
func (h *Palette) Descr() string {
	return "Creates an solid 256x56 image of the color specified."
}

// Re is the regex for matching color messages.
func (h *Palette) Re() string {
	return `(?i)^#[a-f0-9]{6}$`
}

// Match determines if we are asking for a color
func (h *Palette) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *Palette) SetStore(s PluginStore) {}

func (h *Palette) parseHexColor(s string) (*color.RGBA, error) {
	c := &color.RGBA{
		A: 0xff,
	}
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// RespondText to color request events
func (h *Palette) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	const width, height = 256, 56

	img := image.NewRGBA(image.Rect(0, 0, 256, 56))
	color, err := h.parseHexColor(post)
	if err != nil {
		fmt.Println(err)
		SendText(c, ev.RoomID, fmt.Sprintf("%s", err))
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color)
		}
	}

	err = SendImage(c, ev.RoomID, img)
	if err != nil {
		fmt.Println(err)
	}
}

// Name color
func (h *Palette) Name() string {
	return "Palette"
}
