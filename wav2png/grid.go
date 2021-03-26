package wav2png

import (
	"image"
	"image/color"
	"math"
)

func Grid(width, height, padding uint) *image.NRGBA {
	w := width + 2*padding
	h := height + 2*padding
	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	bounds := img.Bounds()

	// fill background (Ref. https://blog.golang.org/image)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 255})
		}
	}

	// vertical lines
	for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(x, y, color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}

	// horizontal lines
	y0 := float64(bounds.Min.Y) + float64(bounds.Max.Y-bounds.Min.Y-1)/2.0
	yt := int(math.Floor(y0))
	yl := int(math.Ceil(y0))

	for y := yt; y > bounds.Min.Y; y -= 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}

	for y := yl; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}

	return img
}
