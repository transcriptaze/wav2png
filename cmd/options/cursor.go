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
	cursor string
	fn     string
}

type CursorFunc func(t, window, duration time.Duration) (time.Duration, float64, float64)

//go:embed cursors/none.png
var cursor_none []byte

//go:embed cursors/green.png
var cursor_green []byte

//go:embed cursors/red.png
var cursor_red []byte

var cursors = map[string][]byte{
	"none":  cursor_none,
	"green": cursor_green,
	"red":   cursor_red,
}

var linear = func(at, window, duration time.Duration) (time.Duration, float64, float64) {
	t := at.Seconds() / duration.Seconds()
	m := 1.0
	c := 0.0
	x := m*t + c

	percentage := at.Seconds() / duration.Seconds()
	offset := at - seconds(percentage*window.Seconds())

	return offset, x, 0.0
}

var centre = func(at, window, duration time.Duration) (time.Duration, float64, float64) {
	t := at.Seconds() / duration.Seconds()
	m := 0.0
	c := 0.5
	x := m*t + c

	percentage := at.Seconds() / duration.Seconds()
	offset := at - seconds(percentage*window.Seconds())

	if at < window/2 {
		percentage := at.Seconds() / window.Seconds()
		shift := 0.5 - percentage

		return offset, x, shift
	}

	if at > (duration - window/2) {
		percentage := (duration - at).Seconds() / window.Seconds()
		shift := -0.5 + percentage

		return offset, x, shift
	}

	return offset, x, 0.0
}

var left = func(at, window, duration time.Duration) (time.Duration, float64, float64) {
	t := at.Seconds() / duration.Seconds()
	m := 0.0
	c := 0.0
	x := m*t + c

	percentage := at.Seconds() / duration.Seconds()
	offset := at - seconds(percentage*window.Seconds())

	if at > (duration - window) {
		percentage := (duration - at).Seconds() / window.Seconds()
		shift := -1.0 + percentage

		return offset, x, shift
	}

	return offset, x, 0.0
}

var right = func(at, window, duration time.Duration) (time.Duration, float64, float64) {
	t := at.Seconds() / duration.Seconds()
	m := 0.0
	c := 1.0
	x := m*t + c

	percentage := at.Seconds() / duration.Seconds()
	offset := at - seconds(percentage*window.Seconds())

	if at < window {
		percentage := at.Seconds() / window.Seconds()
		shift := 1.0 - percentage

		return offset, x, shift
	}

	return offset, x, 0.0
}

var ease = func(at, window, duration time.Duration) (time.Duration, float64, float64) {
	t := at.Seconds() / duration.Seconds()
	m := 0.0
	c := 0.5

	switch {
	case t <= 0.2:
		m = 2.5
		c = 0.0

	case t >= 0.8:
		m = 2.5
		c = -1.5
	}

	x := m*t + c

	percentage := at.Seconds() / duration.Seconds()
	offset := at - seconds(percentage*window.Seconds())

	return offset, x, 0.0
}

func (c Cursor) String() string {
	cursor := "none"
	fn := c.fn

	if c.cursor != "" {
		cursor = c.cursor
	}

	if fn == "" {
		return fmt.Sprintf("%s", cursor)
	} else {
		return fmt.Sprintf("%s:%s", cursor, fn)
	}
}

func (c *Cursor) Set(s string) error {
	tokens := strings.Split(s, ":")

	if len(tokens) > 0 {
		token := tokens[0]
		match := regexp.MustCompile("^(none|green|red)$").FindStringSubmatch(strings.ToLower(token))

		if match != nil && len(match) > 1 {
			c.cursor = match[1]
		} else if info, err := os.Stat(token); os.IsNotExist(err) {
			return fmt.Errorf("Cursor %v does not exist", token)
		} else if info.Mode().IsDir() || !info.Mode().IsRegular() {
			return fmt.Errorf("Cursor file %v is not a file", token)
		} else {
			c.cursor = token
		}
	}

	if len(tokens) > 1 {
		token := tokens[1]
		match := regexp.MustCompile("^(linear|centre|center|left|right|ease)$").FindStringSubmatch(strings.ToLower(token))

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

	case "center":
		return centre

	case "left":
		return left

	case "right":
		return right

	case "ease":
		return ease
	}

	return linear
}

func (c Cursor) Render(h int) *image.NRGBA {
	if b, ok := cursors[c.cursor]; ok {
		return c.make(b, h)
	}

	if b, err := os.ReadFile(c.cursor); err == nil {
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

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
