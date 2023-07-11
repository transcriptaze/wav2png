package styles

import (
	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/wav2png"
)

type Style struct {
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
