package styles

import (
	"image/color"

	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/palettes"
	"github.com/transcriptaze/wav2png/renderers"
	"github.com/transcriptaze/wav2png/renderers/lines"
)

var BLACK = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
var GREEN = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}

type Style struct {
	Name       string
	Width      uint
	Height     uint
	Padding    int
	background string
	grid       string
	renderer   string
}

func NewStyle() Style {
	return Style{
		Name:       "default",
		Width:      800,
		Height:     600,
		Padding:    2,
		background: "black",
		grid:       "square",
		renderer:   "lines, palette:ice, antialias:vertical, vscale:1.0",
	}
}

func (s Style) Fill() fills.FillSpec {
	return fills.NewSolidFill(BLACK)
}

func (s Style) Grid() grids.GridSpec {
	return grids.NewSquareGrid(GREEN, 64, grids.Approximate, false)
}

func (s Style) Renderer() renderers.Renderer {
	return lines.Lines{
		Palette:   palettes.Fire,
		AntiAlias: kernels.Vertical,
		VScale:    1.0,
	}
}
