package renderers

import (
	"image"
)

type Renderer interface {
	Render(audio []float32, width, height, padding int) (*image.NRGBA, error)
}
