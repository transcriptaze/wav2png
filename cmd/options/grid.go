package options

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Grid struct {
	Grid    string `json:"grid"`
	Colour  string `json:"colour"`
	Alpha   uint8  `json:"alpha"`
	Size    string `json:"size"`
	WH      string `json:"wh"`
	Overlay bool   `json:"overlay"`
}

func (g Grid) String() string {
	if g.Grid == "none" {
		return fmt.Sprintf("%v", g.Grid)
	}

	if g.Grid == "square" {
		if g.Overlay {
			return fmt.Sprintf("%v:%v%02x:%v:overlay", g.Grid, g.Colour, g.Alpha, g.Size)
		} else {
			return fmt.Sprintf("%v:%v%02x:%v", g.Grid, g.Colour, g.Alpha, g.Size)
		}
	}

	if g.Grid == "rectangular" {
		if g.Overlay {
			return fmt.Sprintf("%v:%v%02x:%v:overlay", g.Grid, g.Colour, g.Alpha, g.WH)
		} else {
			return fmt.Sprintf("%v:%v%02x:%v", g.Grid, g.Colour, g.Alpha, g.WH)
		}
	}

	return "??"
}

func (g *Grid) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(none|square|rectangular).*").FindStringSubmatch(ss)

	if match != nil && len(match) > 1 {
		switch match[1] {
		case "none":
			g.Grid = "none"

		case "square":
			g.Grid = "square"

			match = regexp.MustCompile("^square:(#[[:xdigit:]]{8}).*").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				color := colour(match[1])
				g.Colour = fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
				g.Alpha = color.A
			}

			match = regexp.MustCompile("^square:#[[:xdigit:]]{8}:([~=≥≤><]?[0-9]+).*").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				fit, size := size(match[1])
				g.Size = fmt.Sprintf("%v%v", fit, size)
			}

			match = regexp.MustCompile("^square:#[[:xdigit:]]{8}:[~=≥≤><]?[0-9]+:(overlay)").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				g.Overlay = true
			}

		case "rectangular":
			g.Grid = "rectangular"

			match = regexp.MustCompile("^rectangular:(#[[:xdigit:]]{8}).*").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				color := colour(match[1])
				g.Colour = fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
				g.Alpha = color.A
			}

			match = regexp.MustCompile("^rectangular:#[[:xdigit:]]{8}:([~=≥≤><]?[0-9]+x[0-9]+).*").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				fit, w, h := wh(match[1])
				g.WH = fmt.Sprintf("%v%vx%v", fit, w, h)
			}

			match = regexp.MustCompile("^rectangular:#[[:xdigit:]]{8}:[~=≥≤><]?[0-9]+x[0-9]+:(overlay)").FindStringSubmatch(ss)
			if match != nil && len(match) > 1 {
				g.Overlay = true
			}

		}
	}

	return nil
}

func (g Grid) GridSpec() wav2png.GridSpec {
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

func size(s string) (wav2png.Fit, uint) {
	fit := wav2png.Approximate
	size := uint(64)

	if matched := regexp.MustCompile(`([~=><≥≤])?\s*([0-9]+)`).FindStringSubmatch(s); matched != nil && len(matched) == 3 {
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

	return fit, size
}

func wh(s string) (wav2png.Fit, uint, uint) {
	fit := wav2png.Approximate
	width := uint(64)
	height := uint(48)

	if matched := regexp.MustCompile(`([~=><≥≤])?\s*([0-9]+)x([0-9]+)`).FindStringSubmatch(s); matched != nil && len(matched) == 4 {
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

	return fit, width, height
}
