package lines

import (
	"image"
	"image/draw"
	"math"
	"time"

	"github.com/transcriptaze/wav2png/wav2png"
)

const (
	RANGE_MIN int32 = -32768
	RANGE_MAX int32 = +32767
	RANGE     int32 = RANGE_MAX - RANGE_MIN + 1
)

type Lines struct {
	Width     int
	Height    int
	Padding   int
	Palette   wav2png.Palette
	FillSpec  wav2png.FillSpec
	GridSpec  wav2png.GridSpec
	AntiAlias wav2png.Kernel
	VScale    float64
}

func (l Lines) Render(samples []float32, fs float64) (*image.NRGBA, error) {
	width := l.Width
	height := l.Height
	padding := l.Padding

	w := width
	h := height
	if padding > 0 {
		w = width - 2*padding
		h = height - 2*padding
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := wav2png.Grid(l.GridSpec, width, height, padding)
	waveform := l.render(samples, fs, w, h)

	x0 := padding
	y0 := padding
	x1 := x0 + w
	y1 := y0 + h

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, y0, x1, y1)
	rectg := img.Bounds()

	wav2png.Fill(img, l.FillSpec)

	if l.GridSpec.Overlay() {
		draw.Draw(img, rect, waveform, origin, draw.Over)
		draw.Draw(img, rectg, grid, origin, draw.Over)
	} else {
		draw.Draw(img, rectg, grid, origin, draw.Over)
		draw.Draw(img, rect, waveform, origin, draw.Over)
	}

	return img, nil
}

func (r Lines) render(samples []float32, fs float64, width, height int) *image.NRGBA {
	volume := r.VScale
	pps := float64(width) / float64(len(samples))
	duration := seconds(float64(len(samples)) / fs)
	l := int(math.Round(math.Ceil(fs*duration.Seconds()) / float64(width)))
	buffer := make([]float32, l)
	waveform := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	colours := r.Palette.Realize()

	x := 0.0
	offset := x / pps
	start := int(math.Round(offset))

	for start < len(samples) {
		end := int(math.Round(offset + 1.0/pps))
		N := copy(buffer, samples[start:end])

		sum := make([]int, height)
		u := vscale(0, -int(height))
		for _, sample := range buffer[0:N] {
			v := int16(32768 * float64(sample) * volume)
			h := vscale(v, -int(height))
			dy := signum(int(h) - int(u))
			for y := int(u); y != int(h); y += dy {
				sum[y]++
			}
		}

		for y := 0; y < height; y++ {
			if sum[y] > 0 {
				l := len(colours)
				i := ceil((l-1)*sum[y], N)
				waveform.Set(int(x+1.0), int(y), colours[i])
			}
		}

		x += 1.0
		offset = x / pps
		start = int(math.Round(offset))
	}

	return wav2png.Antialias(waveform, r.AntiAlias)
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
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
