package wav2png

import (
	"image"
	"image/color"
)

type GridSpec interface {
	Colour() color.NRGBA
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
