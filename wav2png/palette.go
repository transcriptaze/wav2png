package wav2png

import (
	"fmt"
	"image"
	"image/color"
)

type Palette struct {
	colours []color.NRGBA
}

func PaletteFromPng(png image.Image) (*Palette, error) {

	bounds := png.Bounds()
	if bounds.Empty() {
		return nil, fmt.Errorf("cannot create palette from empty PNG")
	}

	h := bounds.Size().Y
	colours := make([]color.NRGBA, h)
	nrgba := color.NRGBAModel

	for i := 0; i < h; i++ {
		colours[i] = nrgba.Convert(png.At(0, i)).(color.NRGBA)
	}

	return &Palette{
		colours: colours,
	}, nil
}

func (p *Palette) realize() []color.NRGBA {
	return p.colours
}

var Mono Palette = Palette{
	colours: []color.NRGBA{
		color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		color.NRGBA{R: 0x80, G: 0x80, B: 0xff, A: 0xff},
	},
}

var Ice Palette = Palette{
	colours: []color.NRGBA{
		color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x02},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x0c},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x16},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x1f},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x29},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x32},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x3c},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x45},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x4f},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x58},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x62},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x6b},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x74},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x7f},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x88},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x91},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0x9b},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xa5},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xae},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xb7},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xc1},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xcb},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xd4},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xde},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xe7},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xf0},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xfa},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
		color.NRGBA{R: 0xc8, G: 0xe8, B: 0xff, A: 0xff},
	},
}
