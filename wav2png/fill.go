package wav2png

import (
	"image"
	"image/color"
)

func Fill(img *image.NRGBA, background color.NRGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, background)
		}
	}
}
