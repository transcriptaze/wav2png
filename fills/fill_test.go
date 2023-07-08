package fills

import (
	"image"
	"image/color"
	"testing"
)

func BenchmarkSolidFill(b *testing.B) {
	img := image.NewNRGBA(image.Rect(0, 0, 2048, 1536))
	fill := SolidFill{
		colour: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255},
	}

	fill.Fill(img)
}
