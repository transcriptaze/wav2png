package options

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"regexp"
	"strings"
	"time"
)

type Cursor struct {
	Cursor string
}

type CursorFunc func(t, offset, window, duration time.Duration) float64

//go:embed cursor_green.png
var green_cursor []byte

//go:embed cursor_red.png
var red_cursor []byte

var cursors = map[string][]byte{
	"green": green_cursor,
	"red":   red_cursor,
}

var linear = func(t, offset, window, duration time.Duration) float64 {
	dt := t - offset
	percentage := dt.Seconds() / window.Seconds()

	return percentage
}

func (c Cursor) String() string {
	return c.Cursor
}

func (c *Cursor) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(none|green|red)$").FindStringSubmatch(ss)

	if match != nil && len(match) > 1 {
		c.Cursor = match[1]
		return nil
	}

	if info, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("Cursor %v does not exist", s)
	} else if info.Mode().IsDir() || !info.Mode().IsRegular() {
		return fmt.Errorf("Cursor file %v is not a file", s)
	} else {
		c.Cursor = s
	}

	return nil
}

func (c Cursor) Fn() CursorFunc {
	return linear
}

func (c Cursor) Render(h int) *image.NRGBA {
	if b, ok := cursors[c.Cursor]; ok {
		return c.make(b, h)
	}

	if b, err := os.ReadFile(c.Cursor); err == nil {
		return c.make(b, h)
	}

	return nil
}

func (c Cursor) make(b []byte, h int) *image.NRGBA {
	cursor, err := png.Decode(bytes.NewBuffer(b))
	if err != nil {
		return nil
	}

	dw := cursor.Bounds().Dx()
	dh := cursor.Bounds().Dy()
	img := image.NewNRGBA(image.Rect(0, 0, dw, h))

	for y := 0; y < h; y += dh {
		draw.Draw(img, image.Rect(0, y, dw, y+dh), cursor, image.Pt(0, 0), draw.Over)
	}

	return img
}
