package options

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
)

type settings struct {
	Size      Size      `json:"size,omitempty"`
	Palette   Palette   `json:"palette,omitempty"`
	Fill      Fill      `json:"fill,omitempty"`
	Padding   Padding   `json:"padding,omitempty"`
	Grid      Grid      `json:"grid,omitempty"`
	Antialias Antialias `json:"antialias,omitempty"`
	Scale     Scale     `json:"scale,omitempty"`
}

type Size struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type Padding int

func (s *settings) Load(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

func colour(s string) color.NRGBA {
	var red uint8
	var green uint8
	var blue uint8
	var alpha uint8

	if _, err := fmt.Sscanf(s, "#%02x%02x%02x%02x", &red, &green, &blue, &alpha); err == nil {
		return color.NRGBA{R: red, G: green, B: blue, A: alpha}
	}

	return color.NRGBA{R: 0, G: 128, B: 0, A: 255}
}