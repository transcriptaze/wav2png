package kernels

import (
	"image"
	"image/color"
)

type Kernel [3][3]uint32

var None = Kernel{
	{0, 0, 0},
	{0, 1, 0},
	{0, 0, 0},
}

var Vertical = Kernel{
	{0, 1, 0},
	{0, 2, 0},
	{0, 1, 0},
}

var Horizontal = Kernel{
	{0, 0, 0},
	{1, 2, 1},
	{0, 0, 0},
}

var Soft = Kernel{
	{1, 2, 1},
	{2, 12, 2},
	{1, 2, 1},
}

func Antialias(img *image.NRGBA, kernel Kernel) *image.NRGBA {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(w, h),
	})

	N := uint32(0)
	for _, row := range kernel {
		for _, k := range row {
			N += k
		}
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r := uint32(0)
			g := uint32(0)
			b := uint32(0)
			a := uint32(0)

			for i, row := range kernel {
				for j, k := range row {
					u := img.NRGBAAt(x+j-1, y+i-1)

					r += k * uint32(u.R)
					g += k * uint32(u.G)
					b += k * uint32(u.B)
					a += k * uint32(u.A)
				}
			}

			out.SetNRGBA(x, y, color.NRGBA{R: uint8(r / N), G: uint8(g / N), B: uint8(b / N), A: uint8(a / N)})
		}
	}

	return out
}
