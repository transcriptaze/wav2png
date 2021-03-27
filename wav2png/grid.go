package wav2png

import (
	"image"
	"image/color"
	"math"
)

func Grid(img *image.NRGBA, colour color.NRGBA) {
	bounds := img.Bounds()

	// vertical lines
	for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(x, y, colour)
		}
	}

	// horizontal lines
	y0 := float64(bounds.Min.Y) + float64(bounds.Max.Y-bounds.Min.Y-1)/2.0
	yt := int(math.Floor(y0))
	yl := int(math.Ceil(y0))

	for y := yt; y > bounds.Min.Y; y -= 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, colour)
		}
	}

	for y := yl; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}
}
