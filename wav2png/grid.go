package wav2png

import (
	"image"
	"image/color"
	"math"
)

type GridSpec interface {
	Colour() color.NRGBA
	Size(bounds image.Rectangle) uint
}

func Grid(img *image.NRGBA, spec GridSpec) {
	bounds := img.Bounds()
	colour := spec.Colour()
	dw := int(spec.Size(bounds))

	// vertical lines
	for x := bounds.Min.X; x < bounds.Max.X; x += dw {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(x, y, colour)
		}
	}

	// horizontal lines
	y0 := float64(bounds.Min.Y) + float64(bounds.Max.Y-bounds.Min.Y-1)/2.0
	yt := int(math.Floor(y0))
	yl := int(math.Ceil(y0))

	for y := yt; y > bounds.Min.Y; y -= dw {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, colour)
		}
	}

	for y := yl; y < bounds.Max.Y; y += dw {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, colour)
		}
	}
}

type SquareGrid struct {
	colour color.NRGBA
	size   uint
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
