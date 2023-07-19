// Package wav2png implements the functions required to render an audio waveform as a PNG image.
// The current implementation supports 16 bit PCM WAV files only.
package wav2png

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/transcriptaze/wav2png/go/kernels"
)

const (
	BITS      int32 = 16
	RANGE_MIN int32 = -32768
	RANGE_MAX int32 = +32767
	RANGE     int32 = RANGE_MAX - RANGE_MIN + 1
)

func Render(pcm []float32, fs float64, width, height int, palette Palette, volume float64) *image.NRGBA {
	pps := float64(width) / float64(len(pcm))
	duration := seconds(float64(len(pcm)) / fs)
	l := int(math.Round(math.Ceil(fs*duration.Seconds()) / float64(width)))
	buffer := make([]float32, l)
	waveform := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	colours := palette.Realize()

	x := 0.0
	offset := x / pps
	start := int(math.Round(offset))

	for start < len(pcm) {
		end := int(math.Round(offset + 1.0/pps))
		N := copy(buffer, pcm[start:end])

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

	return waveform
}

func Antialias(img *image.NRGBA, kernel kernels.Kernel) *image.NRGBA {
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

// rescale converts a PCM value to the 16-bit value used internally.
func rescale(pcmValue int, sampleBitDepth uint) int16 {
	v := int32(pcmValue)
	bits := int32(sampleBitDepth)

	return int16(v * bits / BITS)
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

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
