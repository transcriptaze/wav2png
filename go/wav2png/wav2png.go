// Package wav2png implements the functions required to render an audio waveform as a PNG image.
// The current implementation supports 16 bit PCM WAV files only.
package wav2png

import (
	"errors"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"
)

type Params struct {
	Width   uint
	Height  uint
	Padding uint
}

const (
	BITS      = 17
	RANGE_MIN = -65535
	RANGE_MAX = +65535
)

func Draw(wavfile, pngfile string, params Params) error {
	file, err := os.Open(wavfile)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := wav.NewDecoder(file)

	decoder.FwdToPCM()

	format := decoder.Format()
	bits := decoder.SampleBitDepth()
	duration, err := decoder.Duration()
	if err != nil {
		return err
	}

	fmt.Printf("   File:     %s\n", wavfile)
	fmt.Printf("   PNG:      %s (%d x %d)\n", pngfile, params.Width+2*params.Padding, params.Height+2*params.Padding)
	fmt.Printf("   Format:   %+v\n", *format)
	fmt.Printf("   Bits:     %+v\n", bits)
	fmt.Printf("   Duration: %s\n", duration)
	fmt.Printf("   Length:   %d\n", decoder.PCMLen())

	img, err := plot(decoder, params)
	if err != nil {
		return err
	}

	if img == nil {
		return errors.New("wav2png failed to create image")
	}

	f, err := os.Create(pngfile)
	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, img)
}

func plot(decoder *wav.Decoder, params Params) (*image.NRGBA, error) {
	width := params.Width
	height := params.Height
	padding := params.Padding
	channels := decoder.Format().NumChannels

	bytes := decoder.PCMLen()
	bits := uint(decoder.SampleBitDepth())
	rate := decoder.Format().SampleRate
	samples := bytes / int64(channels*int(bits)/8)
	duration := time.Duration(int64(time.Second) * samples / int64(rate))

	pixels := int64(duration.Round(time.Duration(10)*time.Millisecond)) / 10000000
	msPerPixel := 10 * int64(math.Ceil(float64(pixels)/float64(width)))

	buffer := audio.IntBuffer{Data: make([]int, int64(channels)*441*msPerPixel/10)}
	x := uint(0)

	waveform := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	palette := ice.realize()

	for {
		N, err := decoder.PCMBuffer(&buffer)
		if err != nil {
			return nil, err
		} else if N == 0 {
			break
		}

		sum := make([]int, height)
		u := vscale(rescale(0, bits), height)
		for i := 0; i < N; i += channels {
			v := rescale(buffer.Data[i], bits)
			h := vscale(v, height)
			dy := signum(int(h) - int(u))

			for y := int(u); y != int(h); y += dy {
				sum[y]++
			}
		}

		for y := uint(0); y < height; y++ {
			if sum[y] > 0 {
				l := len(palette)
				i := ceil((l-1)*sum[y], N)
				waveform.Set(int(x), int(height-y-1), palette[i])
			}
		}

		x++
	}

	antialiased := antialias(waveform, soft)
	img := grid(width, height, padding)

	xy := image.Point{0, 0}
	tl := image.Point{int(padding), int(padding)}
	br := image.Point{int(padding + width), int(padding + height)}
	rect := image.Rectangle{tl, br}

	draw.Draw(img, rect, antialiased, xy, draw.Over)

	return img, nil
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

// rescale converts a PCM value to a representative 17-bit integer.
//
// This probably needs some explanation:
//
// The PCM code for a sample represents the 'bottom' value for a bucket of audio values,
// which is all good and fine in practice, but results in some odd quirks when drawing
// waveforms on images with an odd height. For example:
//
// Each PCM code in 16 bit audio sample represents a signal range (scaled to ±1V) of about
// 30µV, which means that '0' is not 0V but is value somewhere in the range 0 - 30µV. Which
// is moot for all practical purposes because the average error of 15µV is entirely lost in
// noise. And is alos irrelevant for PNG files with a height that is an even number because
// the PCM waveform can be correctly rendered as 'two halves', with the '0' X-axis between
// the two halves (wav2png cops out here and renders the 0 X-axis in both halves since it
// doesn't know how to draw between pixels).
//
// For a PNG image with a height that is an odd number however, there is a real pixel that
// is '0', and drawing the PCM '0' code on the '0' X-axis is not correct (not that anybody
// would notice but some of us wake up a 3AM wondering about stuff like this).
//
// Rather than workaround the (whole non-)issue as a bunch of edge cases scattered around
// the code base, wav2png rescales the PCM code to a representative value in the middle of
// the sample bucket with a 17-bit sample depth. 0 is now 0, the sample range is symmetrical
// (±65535), and a 16 bit PCM code becomes the 'next' 17 bit value e.g. PCM code 0x0000 (0)
// becomes 0x00001 (+1) and PCM code 0x0001 (+1) becomes 0x00003 (+3).
//
// As a further example - PCM code 0xffff (-1) represents the 'bucket' -30µV to 0µV, the
// middle of the bucket is -15µV so PCM code 0xffff (-1) is encoded as the 17-bit value
// 0xffff (-1).
//
// As mentioned above, for practical purposes this is essentially irrelevant but it does
// mean the internal arithmetic becomes mentally neat, tidy and and consistent.
func rescale(pcmValue int, sampleBitDepth uint) int32 {
	v := int64(pcmValue) * int64(sampleBitDepth) / (BITS - 1)
	v <<= 1
	v += 1

	return int32(v)
}

// vscale maps the 17-bit internal sample value to a pixel range [0..height). i.e. for a
// height of 256 pixels, vscale maps -65535 to 0 and +65535 to 255.
func vscale(v int32, height uint) int16 {
	vv := int64(v-RANGE_MIN) * int64(height) / ((RANGE_MAX + 1) - (RANGE_MIN - 1))

	return int16(vv)
}

func antialias(img *image.NRGBA, kernel [][]uint32) *image.NRGBA {
	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: w, Y: h},
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
					u := img.At(x+j-1, y+i-1)

					r += k * uint32(u.(color.NRGBA).R)
					g += k * uint32(u.(color.NRGBA).G)
					b += k * uint32(u.(color.NRGBA).B)
					a += k * uint32(u.(color.NRGBA).A)
				}
			}

			out.Set(x, y, color.NRGBA{R: uint8(r / N), G: uint8(g / N), B: uint8(b / N), A: uint8(a / N)})
		}
	}

	return out
}

func grid(width, height, padding uint) *image.NRGBA {
	w := width + 2*padding
	h := height + 2*padding
	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))

	for y := uint(0); y < h; y++ {
		for x := uint(0); x < w; x++ {
			img.Set(int(x), int(y), color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 255})
		}
	}

	for _, y := range []uint{1, 63, 127, 128, 191, 254} {
		for x := uint(0); x < width; x++ {
			img.Set(int(x+padding), int(y+padding), color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}

	for x := uint(0); x <= width; x += 64 {
		for y := uint(1); y < height-1; y++ {
			img.Set(int(x+padding), int(y+padding), color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 255})
		}
	}

	return img
}
