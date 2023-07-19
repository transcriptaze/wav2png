package grids

import (
	"image"
	"image/color"
	"math"
)

type RectangularGrid struct {
	colour  color.NRGBA
	width   uint
	height  uint
	fit     Fit
	overlay bool
}

func NewRectangularGrid(colour color.NRGBA, width, height uint, fit Fit, overlay bool) RectangularGrid {
	return RectangularGrid{
		colour:  colour,
		width:   width,
		height:  height,
		fit:     fit,
		overlay: overlay,
	}
}

func (g RectangularGrid) Colour() color.NRGBA {
	return g.colour
}

func (g RectangularGrid) Overlay() bool {
	return g.overlay
}

func (g RectangularGrid) Border(bounds image.Rectangle, padding int) *image.Rectangle {
	border := image.Rect(bounds.Min.X+padding, bounds.Min.Y+padding, bounds.Max.X-1-padding, bounds.Max.Y-1-padding)

	return &border
}

func (g RectangularGrid) VLines(bounds image.Rectangle, padding int) []int {
	vlines := []int{}

	x0 := bounds.Min.X
	x1 := bounds.Max.X - 1

	border := g.Border(bounds, padding)
	if border != nil {
		x0 = border.Min.X
		x1 = border.Max.X
	}

	N := float64(x1-x0) / float64(g.width)
	dw := float64(x1-x0) / math.Round(N)

	switch g.fit {
	case Exact:
		dw = float64(g.width)

	case AtLeast:
		dw = math.Max(dw, float64(g.width))

	case AtMost:
		dw = math.Min(dw, float64(g.width))

	case LargerThan:
		dw = math.Max(dw, float64(g.width+1))

	case SmallerThan:
		dw = math.Max(dw, float64(g.width-1))
	}

	if dw > 0.0 {
		for line := 1; ; line++ {
			if gx := math.Round(float64(x0) + float64(line)*dw); gx < float64(x1) {
				vlines = append(vlines, int(gx))
				continue
			}

			break
		}
	}

	return vlines
}

func (g RectangularGrid) HLines(bounds image.Rectangle, padding int) []int {
	hlines := []int{}

	y0 := bounds.Min.Y
	y1 := bounds.Max.Y - 1

	border := g.Border(bounds, padding)
	if border != nil {
		y0 = border.Min.Y
		y1 = border.Max.Y
	}

	N := float64(y1-y0) / float64(g.height)
	dw := float64(y1-y0) / math.Round(N)

	switch g.fit {
	case Exact:
		dw = float64(g.height)

	case AtLeast:
		dw = math.Max(dw, float64(g.height))

	case AtMost:
		dw = math.Min(dw, float64(g.height))
	}

	ym := float64(y1-y0+2*padding) / 2.0
	if dw > 0 {
		for line := 0; ; line++ {
			if gy := math.Round(math.Floor(ym) - float64(line)*dw); gy > float64(y0) {
				hlines = append(hlines, int(gy))
				continue
			}

			break
		}

		for line := 0; ; line++ {
			if gy := math.Round(math.Ceil(ym) + float64(line)*dw); gy < float64(y1) {
				hlines = append(hlines, int(gy))
				continue
			}

			break
		}
	}

	return hlines
}
