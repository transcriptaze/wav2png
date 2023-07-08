package columns

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/wav2png"
)

const (
	RANGE_MIN int32 = -32768
	RANGE_MAX int32 = +32767
	RANGE     int32 = RANGE_MAX - RANGE_MIN + 1
)

type Columns struct {
	BarWidth  int
	BarGap    int
	Palette   wav2png.Palette
	AntiAlias kernels.Kernel
	VScale    float64
}

func (c Columns) Render(samples []float32, width, height, padding int) (*image.NRGBA, error) {
	w := width
	h := height
	if padding > 0 {
		w = width - 2*padding
		h = height - 2*padding
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	waveform := c.render(samples, w, h)

	x0 := padding
	y0 := padding
	x1 := x0 + w
	y1 := y0 + h

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, y0, x1, y1)

	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)
	draw.Draw(img, rect, waveform, origin, draw.Over)

	return img, nil
}

func (c Columns) render(samples []float32, width, height int) *image.NRGBA {
	bar := image.NewNRGBA(image.Rect(0, 0, 1, int(height)))
	waveform := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	colours := c.Palette.Realize()
	volume := c.VScale

	scaler := draw.CatmullRom
	column := image.Rect(1, 0, c.BarWidth, height)

	x := 0
	dx := c.BarWidth + c.BarGap
	start := 0

	for start < len(samples) {
		end := (x + dx) * len(samples) / width
		if end > len(samples) {
			end = len(samples)
		}

		sum := make([]int, height)
		u := vscale(0, -int(height))

		for _, sample := range samples[start:end] {
			v := int16(32768 * float64(sample) * volume)
			h := vscale(v, -int(height))
			dy := signum(int(h) - int(u))
			for y := int(u); y != int(h); y += dy {
				sum[y]++
			}
		}

		N := end - start

		draw.Draw(bar, bar.Bounds(), image.Transparent, image.Pt(0, 0), draw.Src)

		for y := 0; y < height; y++ {
			if sum[y] > 0 {
				l := len(colours)
				i := ceil((l-1)*sum[y], N)

				bar.Set(0, y, colours[i])
			}
		}

		xy := image.Pt(x+1, 0)
		scaler.Scale(waveform, column.Add(xy), bar, bar.Bounds(), draw.Over, nil)

		x += dx
		start = end
	}

	return wav2png.Antialias(waveform, c.AntiAlias)
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
func vscale(v int16, height int) int16 {
	h := int32(height)
	vv := int32(v) - RANGE_MIN
	vvv := int16(h * vv / RANGE)

	if height < 0 {
		return vvv - int16(height+1)
	}

	return vvv
}
