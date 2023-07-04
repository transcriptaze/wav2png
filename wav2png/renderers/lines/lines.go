package lines

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"time"

	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/wav2png"
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
	Channels  []int
}

func (l Lines) Render(audio encoding.Audio, from, to time.Duration) (*image.NRGBA, error) {
	w := l.Width
	h := l.Height
	if l.Padding > 0 {
		w = l.Width - 2*l.Padding
		h = l.Height - 2*l.Padding
	}

	fs := audio.SampleRate
	samples := mix(audio, l.Channels...)

	start := int(math.Floor(from.Seconds() * fs))
	if start < 0 || start > len(samples) {
		return nil, fmt.Errorf("start position not in range %v-%v", from, audio.Duration)
	}

	end := int(math.Floor(to.Seconds() * fs))
	if end < 0 || end < start || end > len(samples) {
		return nil, fmt.Errorf("end position not in range %v-%v", from, audio.Duration)
	}

	img := image.NewNRGBA(image.Rect(0, 0, l.Width, l.Height))
	grid := wav2png.Grid(l.GridSpec, l.Width, l.Height, l.Padding)
	waveform := wav2png.Render(samples[start:end], fs, w, h, l.Palette, l.VScale)
	antialiased := wav2png.Antialias(waveform, l.AntiAlias)

	x0 := l.Padding
	y0 := l.Padding
	x1 := x0 + w
	y1 := y0 + h

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, y0, x1, y1)
	rectg := img.Bounds()

	wav2png.Fill(img, l.FillSpec)

	if l.GridSpec.Overlay() {
		draw.Draw(img, rect, antialiased, origin, draw.Over)
		draw.Draw(img, rectg, grid, origin, draw.Over)
	} else {
		draw.Draw(img, rectg, grid, origin, draw.Over)
		draw.Draw(img, rect, antialiased, origin, draw.Over)
	}

	return img, nil
}

func mix(wav encoding.Audio, channels ...int) []float32 {
	L := wav.Length
	N := float64(len(channels))
	samples := make([]float32, L)

	if len(wav.Samples) < 2 {
		return wav.Samples[0]
	}

	for i := 0; i < L; i++ {
		sample := 0.0
		for _, ch := range channels {
			sample += float64(wav.Samples[ch-1][i])
		}

		samples[i] = float32(sample / N)
	}

	return samples
}
