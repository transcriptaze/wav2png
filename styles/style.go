package styles

import (
	"image/color"

	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/wav2png"
)

var BLACK = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}

type Style struct {
	Name       string
	Width      uint
	Height     uint
	Padding    int
	Background fills.FillSpec
	Grid       grids.GridSpec
}

type LinesStyle struct {
	Style
	Palette   wav2png.Palette
	Antialias kernels.Kernel
	VScale    float64
}

func NewStyle() Style {
	return Style{
		Name:       "default",
		Width:      800,
		Height:     600,
		Padding:    2,
		Background: fills.NewSolidFill(BLACK),
	}
}
