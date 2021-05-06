package wav2png

import (
	"image"
	"image/color"
)

type NoGrid struct {
	padding int
}

func NewNoGrid(padding int) GridSpec {
	return NoGrid{
		padding: padding,
	}
}

func (g NoGrid) Colour() color.NRGBA {
	return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
}

func (g NoGrid) Padding() int {
	return g.padding
}

func (g NoGrid) Overlay() bool {
	return false
}

func (g NoGrid) Border(bounds image.Rectangle) *image.Rectangle {
	return nil
}

func (g NoGrid) VLines(bounds image.Rectangle) []int {
	return []int{}
}

func (g NoGrid) HLines(bounds image.Rectangle) []int {
	return []int{}
}
