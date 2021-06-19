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

type Fill struct {
	Fill   string `json:"fill"`
	Colour string `json:"colour"`
	Alpha  uint8  `json:"alpha"`
}

type Antialias struct {
	Type   string `json:"type"`
	Kernel wav2png.Kernel
}

type Scale struct {
	Horizontal float64 `json:"horizontal"`
	Vertical   float64 `json:"vertical"`
}

func (f *Fill) FillSpec() wav2png.FillSpec {
	colour := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}

	red := uint8(0)
	green := uint8(0)
	blue := uint8(0)
	alpha := f.Alpha
	if _, err := fmt.Sscanf(f.Colour, "#%02x%02x%02x", &red, &green, &blue); err == nil {
		colour = color.NRGBA{R: red, G: green, B: blue, A: alpha}
	}

	return wav2png.NewSolidFill(colour)
}
