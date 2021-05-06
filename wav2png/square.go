package wav2png

import (
	"image"
	"image/color"
	"math"
)

type SquareGrid struct {
	colour  color.NRGBA
	padding int
	size    uint
	fit     Fit
	overlay bool
}

type Fit int

const (
	Approximate Fit = iota
	Exact
	AtLeast
	AtMost
)

func NewSquareGrid(colour color.NRGBA, size uint, padding int, fit Fit, overlay bool) GridSpec {
	return SquareGrid{
		colour:  colour,
		padding: padding,
		size:    size,
		fit:     fit,
		overlay: overlay,
	}
}

func (g SquareGrid) Colour() color.NRGBA {
	return g.colour
}

func (g SquareGrid) Padding() int {
	return g.padding
}

func (g SquareGrid) Overlay() bool {
	return g.overlay
}

func (g SquareGrid) Border(bounds image.Rectangle) *image.Rectangle {
	padding := g.padding
	border := image.Rect(bounds.Min.X+padding, bounds.Min.Y+padding, bounds.Max.X-1-padding, bounds.Max.Y-1-padding)

	return &border
}

func (g SquareGrid) VLines(bounds image.Rectangle) []int {
	vlines := []int{}

	x0 := bounds.Min.X
	x1 := bounds.Max.X - 1

	border := g.Border(bounds)
	if border != nil {
		x0 = border.Min.X
		x1 = border.Max.X
	}

	N := float64(x1-x0) / float64(g.size)
	dw := float64(x1-x0) / math.Round(N)

	switch g.fit {
	case Exact:
		dw = float64(g.size)

	case AtLeast:
		dw = math.Max(dw, float64(g.size))

	case AtMost:
		dw = math.Min(dw, float64(g.size))
	}

	if dw > 0.0 {
		for line := 1; ; line++ {
			if gx := math.Round(float64(x0) + float64(line)*dw); gx < float64(x1) {
				vlines = append(vlines, int(gx))
				continue
			}

			break
		}
	}

	return vlines
}

func (g SquareGrid) HLines(bounds image.Rectangle) []int {
	hlines := []int{}

	y0 := bounds.Min.Y
	y1 := bounds.Max.Y - 1

	padding := g.padding
	border := g.Border(bounds)
	if border != nil {
		y0 = border.Min.Y
		y1 = border.Max.Y
	}

	N := float64(y1-y0) / float64(g.size)
	dw := float64(y1-y0) / math.Round(N)

	switch g.fit {
	case Exact:
		dw = float64(g.size)

	case AtLeast:
		dw = math.Max(dw, float64(g.size))

	case AtMost:
		dw = math.Min(dw, float64(g.size))
	}

	ym := float64(y1-y0+2*padding) / 2.0
	if dw > 0 {
		for line := 0; ; line++ {
			if gy := math.Round(math.Floor(ym) - float64(line)*dw); gy > float64(y0) {
				hlines = append(hlines, int(gy))
				continue
			}

			break
		}

		for line := 0; ; line++ {
			if gy := math.Round(math.Ceil(ym) + float64(line)*dw); gy < float64(y1) {
				hlines = append(hlines, int(gy))
				continue
			}

			break
		}
	}

	return hlines
}
