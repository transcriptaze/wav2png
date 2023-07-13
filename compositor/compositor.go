package compositor

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/renderers"
	"github.com/transcriptaze/wav2png/styles"
)

type Compositor struct {
	width      uint
	height     uint
	padding    int
	background fills.FillSpec
	grid       grids.GridSpec
	renderer   renderers.Renderer
}

func FromStyle(style styles.Style) Compositor {
	return Compositor{
		width:      style.Width,
		height:     style.Height,
		padding:    style.Padding,
		background: style.Fill(),
		grid:       style.Grid(),
		renderer:   style.Renderer(),
	}
}

func NewCompositor(width uint, height uint, padding int, background fills.FillSpec, grid grids.GridSpec, renderer renderers.Renderer) Compositor {
	return Compositor{
		width:      width,
		height:     height,
		padding:    padding,
		background: background,
		grid:       grid,
		renderer:   renderer,
	}
}

func (c Compositor) Render(samples []float32) (*image.NRGBA, error) {
	width := int(c.width)
	height := int(c.height)
	padding := c.padding

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := grids.Grid(c.grid, width, height, padding)

	if waveform, err := c.renderer.Render(samples, width, height, padding); err != nil {
		return nil, err
	} else {
		origin := image.Pt(0, 0)
		bounds := img.Bounds()

		fills.Fill(img, c.background)

		if c.grid.Overlay() {
			draw.Draw(img, bounds, waveform, origin, draw.Over)
			draw.Draw(img, bounds, grid, origin, draw.Over)
		} else {
			draw.Draw(img, bounds, grid, origin, draw.Over)
			draw.Draw(img, bounds, waveform, origin, draw.Over)
		}

		return img, nil
	}
}
