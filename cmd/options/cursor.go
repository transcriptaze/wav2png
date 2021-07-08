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
	fn     string
}

type CursorFunc func(t, offset, window, duration time.Duration) (float64, float64)

//go:embed cursor_green.png
var green_cursor []byte

//go:embed cursor_red.png
var red_cursor []byte

var cursors = map[string][]byte{
	"green": green_cursor,
	"red":   red_cursor,
}

var linear = func(t, offset, window, duration time.Duration) (float64, float64) {
	dt := t - offset
	percentage := dt.Seconds() / window.Seconds()

	return percentage, 0.0
}

var centre = func(t, offset, window, duration time.Duration) (float64, float64) {
	if t < window/2 {
		percentage := t.Seconds() / window.Seconds()
		shift := 0.5 - percentage

		return 0.5, shift
	}

	if t > (duration - window/2) {
		percentage := (duration - t).Seconds() / window.Seconds()
		shift := -0.5 + percentage

		return 0.5, shift
	}

	return 0.5, 0.0
}

var left = func(t, offset, window, duration time.Duration) (float64, float64) {
	if t > (duration - window) {
		percentage := (duration - t).Seconds() / window.Seconds()
		shift := -1.0 + percentage

		return 0.0, shift
	}

	return 0.0, 0.0
}

func (c Cursor) String() string {
	return c.Cursor
}

func (c *Cursor) Set(s string) error {
	tokens := strings.Split(s, ":")

	if len(tokens) > 0 {
		token := tokens[0]
		match := regexp.MustCompile("^(none|green|red)$").FindStringSubmatch(strings.ToLower(token))

		if match != nil && len(match) > 1 {
			c.Cursor = match[1]
		} else if info, err := os.Stat(token); os.IsNotExist(err) {
			return fmt.Errorf("Cursor %v does not exist", token)
		} else if info.Mode().IsDir() || !info.Mode().IsRegular() {
			return fmt.Errorf("Cursor file %v is not a file", token)
		} else {
			c.Cursor = token
		}
	}

	if len(tokens) > 1 {
		token := tokens[1]
		match := regexp.MustCompile("^(linear|centre|center|left|right)$").FindStringSubmatch(strings.ToLower(token))

		if match != nil && len(match) > 1 {
			c.fn = match[1]
		}
	}

	return nil
}

func (c Cursor) Fn() CursorFunc {
	switch c.fn {
	case "centre":
		return centre

	case "left":
		return left
	}

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
