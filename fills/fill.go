package fills

import (
	"image"
	"image/color"
)

type FillSpec interface {
	Fill(img *image.NRGBA)
}

func Fill(img *image.NRGBA, spec FillSpec) {
	spec.Fill(img)
}

type SolidFill struct {
	colour color.NRGBA
}

func NewSolidFill(colour color.NRGBA) SolidFill {
	return SolidFill{
		colour: colour,
	}
}

func (f SolidFill) Fill(img *image.NRGBA) {
	bounds := img.Bounds()
	background := f.colour

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, background)
		}
	}
}
