package wav2png

import (
	"image"
	"image/color"
	"math"
)

type GridSpec interface {
	Colour() color.NRGBA
	Size(bounds image.Rectangle) uint
	Padding(bounds image.Rectangle) int
	Border(bounds image.Rectangle) *image.Rectangle
	VLines(bounds image.Rectangle) []int
	HLines(bounds image.Rectangle) []int
}

func Grid(img *image.NRGBA, spec GridSpec) {
	bounds := img.Bounds()
	colour := spec.Colour()

	// calculate grid metrics
	x0 := bounds.Min.X
	x1 := bounds.Max.X - 1

	y0 := bounds.Min.Y
	y1 := bounds.Max.Y - 1

	border := spec.Border(bounds)
	if border != nil {
		x0 = border.Min.X
		x1 = border.Max.X

		y0 = border.Min.Y
		y1 = border.Max.Y
	}

	// vertical lines
	vlines := spec.VLines(bounds)
	for _, x := range vlines {
		for y := y0; y <= y1; y++ {
			img.Set(x, y, colour)
		}
	}

	// horizontal lines
	// c := color.NRGBA{R: 0xff, G: 0x80, B: 0x00, A: 0xff}
	hlines := spec.HLines(bounds)
	for _, y := range hlines {
		for x := x0; x <= x1; x++ {
			img.Set(x, y, colour)
		}
	}

	// border
	if border != nil {
		for x := border.Min.X; x <= border.Max.X; x++ {
			img.Set(x, border.Min.Y, colour)
			img.Set(x, border.Max.Y, colour)
		}

		for y := border.Min.Y; y <= border.Max.Y; y++ {
			img.Set(border.Min.X, y, colour)
			img.Set(border.Max.X, y, colour)
		}
	}

}

type SquareGrid struct {
	colour  color.NRGBA
	size    uint
	padding int
}

func NewSquareGrid(colour color.NRGBA, size uint, padding int) SquareGrid {
	return SquareGrid{
		colour:  colour,
		size:    size,
		padding: padding,
	}
}

func (g SquareGrid) Colour() color.NRGBA {
	return g.colour
}

func (g SquareGrid) Size(bound image.Rectangle) uint {
	return g.size
}

func (g SquareGrid) Padding(bound image.Rectangle) int {
	return g.padding
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

	padding := g.padding
	border := g.Border(bounds)
	if border != nil {
		x0 = border.Min.X
		x1 = border.Max.X
	}

	if dw := (x1 - x0) / 10.0; dw > 0 {
		for gx := x0 + padding + dw; gx < x1; gx += dw {
			vlines = append(vlines, gx)
		}
	}

	return vlines
}

func (g SquareGrid) HLines(bounds image.Rectangle) []int {
	hlines := []int{}

	x0 := bounds.Min.X
	x1 := bounds.Max.X - 1

	y0 := bounds.Min.Y
	y1 := bounds.Max.Y - 1

	padding := g.padding
	border := g.Border(bounds)
	if border != nil {
		x0 = border.Min.X
		x1 = border.Max.X

		y0 = border.Min.Y
		y1 = border.Max.Y
	}

	ym := float64(y1-y0+2*padding) / 2.0
	if dw := (x1 - x0) / 10.0; dw > 0 {
		y := int(math.Floor(ym))
		for gy := y; gy > y0; gy -= dw {
			hlines = append(hlines, gy)
		}

		y = int(math.Ceil(ym))
		for gy := y; gy < y1; gy += dw {
			hlines = append(hlines, gy)
		}
	}

	return hlines
}
