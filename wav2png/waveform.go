package wav2png

import (
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

func Plot(wavfile, pngfile string, width, height, padding uint) error {
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
	fmt.Printf("   PNG:      %s (%d x %d)\n", pngfile, width+2*padding, height+2*padding)
	fmt.Printf("   Format:   %+v\n", *format)
	fmt.Printf("   Bits:     %+v\n", bits)
	fmt.Printf("   Duration: %s\n", duration)
	fmt.Printf("   Length:   %d\n", decoder.PCMLen())

	return plot(decoder, pngfile, width, height, padding)
}

func plot(decoder *wav.Decoder, pngfile string, width, height, padding uint) error {
	w := width + 2*padding
	h := height + 2*padding

	bytes := decoder.PCMLen()
	bits := decoder.SampleBitDepth()
	channels := decoder.Format().NumChannels
	rate := decoder.Format().SampleRate
	samples := bytes / int64(channels*int(bits)/8)
	duration := time.Duration(1000000000 * samples / int64(rate))

	pixels := int64(duration.Round(time.Duration(10)*time.Millisecond)) / 10000000
	msPerPixel := 10 * int64(math.Ceil(float64(pixels)/float64(width)))

	buffer := audio.IntBuffer{Data: make([]int, int64(channels)*441*msPerPixel/10)}
	x := uint(0)

	waveform := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	palette := ice.realize()

	for {
		N, err := decoder.PCMBuffer(&buffer)
		if err != nil {
			return err
		} else if N == 0 {
			break
		}

		sum := make([]int, height)
		u := vscale(0, height)
		for i := 0; i < N; i += channels {
			v := vscale(buffer.Data[i], height)
			dy := signum(v - u)

			for y := u; y != v; y += dy {
				sum[y]++
			}
		}

		pixels := make([]color.NRGBA, height+2)
		px := 0
		pixels[px] = palette[0]
		px++
		for y := uint(0); y < height; y++ {
			pixels[px] = palette[0]
			if sum[y] > 0 {
				l := len(palette)
				i := ceil((l-1)*sum[y], N)
				pixels[px] = palette[i]
			}
			px++
		}
		pixels[px] = palette[0]

		p := antialias(pixels)

		for y := uint(0); y < height; y++ {
			waveform.Set(int(x+padding), int(height-y-1+padding), p[y+1])
		}

		x++
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))

	grid(img, width, height, padding)
	draw.Draw(img, waveform.Bounds(), waveform, waveform.Bounds().Min, draw.Over)

	f, err := os.Create(pngfile)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
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

func vscale(v int, height uint) int {
	return int(height) * (v + 32768) / 65536
}

func antialias(pixels []color.NRGBA) []color.NRGBA {
	p := make([]color.NRGBA, len(pixels))
	L := len(pixels) - 1
	i := 0

	p[i] = pixels[i]
	i++
	for ; i < L; i++ {
		r := (2*uint32(pixels[i].R) + uint32(pixels[i-1].R) + uint32(pixels[i+1].R)) / 4
		g := (2*uint32(pixels[i].G) + uint32(pixels[i-1].G) + uint32(pixels[i+1].G)) / 4
		b := (2*uint32(pixels[i].B) + uint32(pixels[i-1].B) + uint32(pixels[i+1].B)) / 4
		a := (2*uint32(pixels[i].A) + uint32(pixels[i-1].A) + uint32(pixels[i+1].A)) / 4
		p[i] = color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}
	p[i] = pixels[i]

	return p
}

func grid(img *image.NRGBA, width, height, padding uint) {
	w := width + 2*padding
	h := height + 2*padding

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

}
