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
	return "Creates an solid 56x56 image of the color specified."
}

// Re is the regex for matching color messages.
func (h *Palette) Re() string {
	return `(?i)^#[a-f0-9]{6}$`
}

// Match determines if we are asking for a color
func (h *Palette) Match(_, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *Palette) SetStore(_ PluginStore) {}

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

func isEdge(x, y int) bool {
	if x == 0 || x == 55 {
		return true
	}

	if y == 0 || y == 55 {
		return true
	}

	return false
}

// RespondText to color request events
func (h *Palette) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	const width, height = 56, 56

	img := image.NewRGBA(image.Rect(0, 0, 56, 56))
	border := &color.RGBA{
		R: 0x00,
		G: 0x00,
		B: 0x00,
		A: 0xff,
	}
	clr, err := h.parseHexColor(post)
	if err != nil {
		fmt.Println(err)
		return SendText(c, ev.RoomID, fmt.Sprintf("%s", err))
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isEdge(x, y) {
				img.Set(x, y, border)
			} else {
				img.Set(x, y, clr)
			}
		}
	}

	err = SendImage(c, ev.RoomID, img)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Name color
func (h *Palette) Name() string {
	return "Palette"
}
