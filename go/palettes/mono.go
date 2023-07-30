package palettes

import (
	"image/color"
)

var Mono Palette = Palette{
	colours: []color.NRGBA{
		color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		color.NRGBA{R: 0x80, G: 0x80, B: 0xff, A: 0xff},
	},
}
