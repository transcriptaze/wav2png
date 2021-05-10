package wav2png

import (
	"image"
	"image/color"
)

type NoGrid struct {
}

func NewNoGrid() GridSpec {
	return NoGrid{}
}

func (g NoGrid) Colour() color.NRGBA {
	return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
}

func (g NoGrid) Overlay() bool {
	return false
}

func (g NoGrid) Border(bounds image.Rectangle, padding int) *image.Rectangle {
	return nil
}

func (g NoGrid) VLines(bounds image.Rectangle, padding int) []int {
	return []int{}
}

func (g NoGrid) HLines(bounds image.Rectangle, padding int) []int {
	return []int{}
}
