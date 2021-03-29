package wav2png

import (
	"image"
	"image/color"
)

type FillSpec interface {
	Colour() color.NRGBA
}

func Fill(img *image.NRGBA, spec FillSpec) {
	bounds := img.Bounds()
	background := spec.Colour()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, background)
		}
	}
}

type SolidFill struct {
	colour color.NRGBA
}

func NewSolidFill(colour color.NRGBA) SolidFill {
	return SolidFill{
		colour: colour,
	}
}

func (f SolidFill) Colour() color.NRGBA {
	return f.colour
}
