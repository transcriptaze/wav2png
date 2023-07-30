package lines

import (
	"golang.org/x/image/draw"
	"image"

	"github.com/transcriptaze/wav2png/go/kernels"
	"github.com/transcriptaze/wav2png/go/palettes"
)

const (
	RANGE_MIN int32 = -32768
	RANGE_MAX int32 = +32767
	RANGE     int32 = RANGE_MAX - RANGE_MIN + 1
)

type Lines struct {
	Palette   palettes.Palette
	AntiAlias kernels.Kernel
}

func (l Lines) Render(samples []float32, width, height, padding int, vscale float64) (*image.NRGBA, error) {
	w := width
	h := height
	if padding > 0 {
		w = width - 2*padding
		h = height - 2*padding
	}

	x0 := padding
	y0 := padding
	x1 := x0 + w
	y1 := y0 + h

	rect := image.Rect(x0, y0, x1, y1)

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	waveform := l.render(img.SubImage(rect).(*image.NRGBA), samples, vscale)

	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)
	draw.Draw(img, rect, waveform, image.Pt(0, 0), draw.Over)

	return img, nil
}

func (r Lines) render(waveform *image.NRGBA, samples []float32, vscale float64) *image.NRGBA {
	width := waveform.Bounds().Dx()
	height := waveform.Bounds().Dy()
	colours := r.Palette.Realize()

	x := 0
	dx := 1
	start := 0

	for start < len(samples) {
		end := (x + dx) * len(samples) / width

		sum := make([]int, height)
		u := scale(0, -int(height))
		for _, sample := range samples[start:end] {
			v := int16(32768 * float64(sample) * vscale)
			h := scale(v, -int(height))
			dy := signum(int(h) - int(u))
			for y := int(u); y != int(h); y += dy {
				sum[y]++
			}
		}

		N := end - start
		for y := 0; y < height; y++ {
			if sum[y] > 0 {
				l := len(colours)
				i := ceil((l-1)*sum[y], N)
				waveform.Set(x+dx, y, colours[i])
			}
		}

		x += dx
		start = end
	}

	return kernels.Antialias(waveform, r.AntiAlias)
}

func signum(N int) int {
	if N < 0 {
		return -1
	}

	return +1
}

func ceil(p int, q int) int {
	d := p / q
	r := p % q

	if r > 0 {
		return d + 1
	}

	return d
}

// vscale maps the 16-bit internal sample value to a pixel range [0..height). i.e. for a
// height of 256 pixels, vscale maps -32768 to 0 and +32767 to 255. A negative height
// 'flips' the conversion e.g. for height of -256, -32768 is mapped to 255 and +32767 is
// mapped to 0.
func scale(v int16, height int) int16 {
	h := int32(height)
	vv := int32(v) - RANGE_MIN
	vvv := int16(h * vv / RANGE)

	if height < 0 {
		return vvv - int16(height+1)
	}

	return vvv
}
