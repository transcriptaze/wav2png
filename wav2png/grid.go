package wav2png

import (
	"image"
	"image/color"
	"math"
)

type GridSpec interface {
	Colour() color.NRGBA
	Size(bounds image.Rectangle) uint
	Padding(bounds image.Rectangle) uint
}

func Grid(img *image.NRGBA, spec GridSpec) {
	bounds := img.Bounds()
	colour := spec.Colour()
	dw := int(spec.Size(bounds))
	padding := int(spec.Padding(bounds))

	// vertical lines
	for x := bounds.Min.X + padding; x < bounds.Max.X-padding; x += dw {
		for y := bounds.Min.Y + padding; y < bounds.Max.Y-padding; y++ {
			img.Set(x, y, colour)
		}
	}

	// horizontal lines
	y0 := float64(bounds.Min.Y) + float64(bounds.Max.Y-bounds.Min.Y-1)/2.0
	yt := int(math.Floor(y0))
	yl := int(math.Ceil(y0))

	for y := yt; y > bounds.Min.Y+padding; y -= dw {
		for x := bounds.Min.X + padding; x < bounds.Max.X-padding; x++ {
			img.Set(x, y, colour)
		}
	}

	for y := yl; y < bounds.Max.Y-padding; y += dw {
		for x := bounds.Min.X + padding; x < bounds.Max.X-padding; x++ {
			img.Set(x, y, colour)
		}
	}
}

type SquareGrid struct {
	colour  color.NRGBA
	size    uint
	padding uint
}

func NewSquareGrid(colour color.NRGBA, size uint) SquareGrid {
	return SquareGrid{
		colour: colour,
		size:   size,
	}
}

func (g SquareGrid) Colour() color.NRGBA {
	return g.colour
}

func (g SquareGrid) Size(bound image.Rectangle) uint {
	return g.size
}

func (g SquareGrid) Padding(bound image.Rectangle) uint {
	return g.padding
}
