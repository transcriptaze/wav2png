package columns

import (
	"image"
	"image/draw"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Columns struct {
	Width     int
	Height    int
	Padding   int
	Palette   wav2png.Palette
	FillSpec  wav2png.FillSpec
	GridSpec  wav2png.GridSpec
	AntiAlias wav2png.Kernel
	VScale    float64
}

func (c Columns) Render(samples []float32, fs float64) (*image.NRGBA, error) {
	width := c.Width
	height := c.Height
	padding := c.Padding

	w := width
	h := height
	if padding > 0 {
		w = width - 2*padding
		h = height - 2*padding
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := wav2png.Grid(c.GridSpec, width, height, padding)
	waveform := c.render(samples, fs, w, h)

	x0 := padding
	y0 := padding
	x1 := x0 + w
	y1 := y0 + h

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, y0, x1, y1)
	rectg := img.Bounds()

	wav2png.Fill(img, c.FillSpec)

	if c.GridSpec.Overlay() {
		draw.Draw(img, rect, waveform, origin, draw.Over)
		draw.Draw(img, rectg, grid, origin, draw.Over)
	} else {
		draw.Draw(img, rectg, grid, origin, draw.Over)
		draw.Draw(img, rect, waveform, origin, draw.Over)
	}

	return img, nil
}

func (c Columns) render(samples []float32, fs float64, w, h int) *image.NRGBA {
	waveform := wav2png.Render(samples, fs, w, h, c.Palette, c.VScale)

	return wav2png.Antialias(waveform, c.AntiAlias)
}
