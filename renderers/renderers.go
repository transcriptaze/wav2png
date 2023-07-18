package renderers

import (
	"image"
)

type Renderer interface {
	Render(audio []float32, width, height, padding int, scale float64) (*image.NRGBA, error)
}
