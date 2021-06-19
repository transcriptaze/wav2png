package options

import (
	"fmt"
	"image/color"

	"github.com/transcriptaze/wav2png/wav2png"
)

type settings struct {
	Size Size `json:"size"`
	// Palettes   Palettes  `json:"palettes"`
	Fill      Fill      `json:"fill"`
	Padding   Padding   `json:"padding"`
	Grid      Grid      `json:"grid"`
	Antialias Antialias `json:"antialias"`
	Scale     Scale     `json:"scale"`
}

type Size struct {
	width  int
	height int
}

type Padding int

type Antialias struct {
	Type   string `json:"type"`
	Kernel wav2png.Kernel
}

type Scale struct {
	Horizontal float64 `json:"horizontal"`
	Vertical   float64 `json:"vertical"`
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
