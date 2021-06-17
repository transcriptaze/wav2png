package main

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Settings struct {
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

type Grid struct {
	Grid    string `json:"grid"`
	Colour  string `json:"colour"`
	Alpha   uint8  `json:"alpha"`
	Size    string `json:"size"`
	WH      string `json:"wh"`
	Overlay bool   `json:"overlay"`
}

type Antialias struct {
	Type   string `json:"type"`
	kernel wav2png.Kernel
}

type Scale struct {
	Horizontal float64 `json:"horizontal"`
	Vertical   float64 `json:"vertical"`
}

func (f *Fill) fillspec() wav2png.FillSpec {
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

func (g *Grid) gridspec() wav2png.GridSpec {
	// ... overlay
	overlay := g.Overlay

	// ... colour
	red := uint8(0)
	green := uint8(128)
	blue := uint8(0)
	alpha := g.Alpha

	colour := color.NRGBA{R: red, G: green, B: blue, A: alpha}

	if _, err := fmt.Sscanf(g.Colour, "#%02x%02x%02x", &red, &green, &blue); err == nil {
		colour = color.NRGBA{R: red, G: green, B: blue, A: g.Alpha}
	}

	// ... no grid
	if g.Grid == "none" {
		return wav2png.NewNoGrid()
	}

	// ... rectangular
	if g.Grid == "rectangular" {
		fit := wav2png.Approximate
		width := uint(64)
		height := uint(48)

		if matched := regexp.MustCompile(`([~=><≥≤])?\s*([0-9]+)x([0-9]+)`).FindStringSubmatch(g.WH); matched != nil && len(matched) == 4 {
			switch matched[1] {
			case "~":
				fit = wav2png.Approximate
			case "=":
				fit = wav2png.Exact
			case "≥":
				fit = wav2png.AtLeast
			case "≤":
				fit = wav2png.AtMost
			case ">":
				fit = wav2png.LargerThan
			case "<":
				fit = wav2png.SmallerThan
			}

			if v, err := strconv.ParseUint(matched[2], 10, 32); err == nil && v >= 16 && v <= 1024 {
				width = uint(v)
			}

			if v, err := strconv.ParseUint(matched[3], 10, 32); err == nil && v >= 16 && v <= 1024 {
				height = uint(v)
			}
		}

		return wav2png.NewRectangularGrid(colour, width, height, fit, overlay)
	}

	// ... default to square
	fit := wav2png.Approximate
	size := uint(64)

	if matched := regexp.MustCompile(`([~=><≥≤])?\s*([0-9]+)`).FindStringSubmatch(g.Size); matched != nil && len(matched) == 3 {
		switch matched[1] {
		case "~":
			fit = wav2png.Approximate
		case "=":
			fit = wav2png.Exact
		case "≥":
			fit = wav2png.AtLeast
		case "≤":
			fit = wav2png.AtMost
		case ">":
			fit = wav2png.LargerThan
		case "<":
			fit = wav2png.SmallerThan
		}

		if v, err := strconv.ParseUint(matched[2], 10, 32); err == nil && v >= 16 && v <= 1024 {
			size = uint(v)
		}
	}

	return wav2png.NewSquareGrid(colour, size, fit, overlay)
}
