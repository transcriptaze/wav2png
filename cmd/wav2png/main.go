package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/transcriptaze/wav2png/cmd/wav2png/options"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/wav2png"
)

type audio struct {
	sampleRate float64
	format     string
	channels   int
	duration   time.Duration
	length     int
	samples    [][]float32
}

func main() {
	options := options.NewOptions()
	if err := options.Parse(); err != nil {
		usage()
		os.Exit(1)
	}

	w, err := read(options.WAV)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	} else if w == nil {
		fmt.Printf("\n   ERROR: unable to read WAV file\n")
		os.Exit(1)
	}

	if options.Debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", options.WAV)
		fmt.Printf("   Channels:    %v\n", w.channels)
		fmt.Printf("   Format:      %v\n", w.format)
		fmt.Printf("   Sample Rate: %v\n", w.sampleRate)
		fmt.Printf("   Duration:    %v\n", w.duration)
		fmt.Printf("   Samples:     %v\n", w.length)
		fmt.Printf("   PNG:         %v\n", options.PNG)
		fmt.Println()
	}

	img, err := render(*w, options)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	if err := write(img, options.PNG); err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}
}

func render(wav audio, options options.Options) (*image.NRGBA, error) {
	width := int(options.Width)
	height := int(options.Height)
	padding := options.Padding
	fillspec := options.FillSpec
	gridspec := options.GridSpec
	kernel := options.Antialias
	vscale := options.VScale
	palette := options.Palette

	w := width
	h := height
	if padding > 0 {
		w = width - padding
		h = height - padding
	}

	fs := wav.sampleRate
	samples := mix(wav, options.Mix.Channels()...)
	start := 0
	end := len(samples)

	if options.From != nil {
		v := int(math.Floor(options.From.Seconds() * fs))
		if v > 0 && v <= len(samples) {
			start = v
		}
	}

	if options.To != nil {
		v := int(math.Floor(options.To.Seconds() * fs))
		if v < start {
			end = start
		} else if v <= len(samples) {
			end = v
		}
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := wav2png.Grid(gridspec, width, height, padding)
	waveform := wav2png.Render(samples[start:end], fs, w, h, palette, vscale)
	antialiased := wav2png.Antialias(waveform, kernel)

	origin := image.Pt(0, 0)
	rect := image.Rect(padding, padding, w, h)
	rectg := img.Bounds()

	wav2png.Fill(img, fillspec)

	if gridspec.Overlay() {
		draw.Draw(img, rect, antialiased, origin, draw.Over)
		draw.Draw(img, rectg, grid, origin, draw.Over)
	} else {
		draw.Draw(img, rectg, grid, origin, draw.Over)
		draw.Draw(img, rect, antialiased, origin, draw.Over)
	}

	return img, nil
}

func read(wavfile string) (*audio, error) {
	file, err := os.Open(wavfile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	w, err := wav.Decode(file)
	if err != nil {
		return nil, err
	}

	return &audio{
		sampleRate: float64(w.Format.SampleRate),
		format:     fmt.Sprintf("%v", w.Format),
		channels:   int(w.Format.Channels),
		duration:   w.Duration(),
		length:     w.Frames(),
		samples:    w.Samples,
	}, nil
}

func write(img *image.NRGBA, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, img)
}

func mix(wav audio, channels ...int) []float32 {
	L := wav.length
	N := float64(len(channels))
	samples := make([]float32, L)

	if len(wav.samples) < 2 {
		return wav.samples[0]
	}

	for i := 0; i < L; i++ {
		sample := 0.0
		for _, ch := range channels {
			sample += float64(wav.samples[ch-1][i])
		}

		samples[i] = float32(sample / N)
	}

	return samples
}

func usage() {
	println()
	println("   Usage: waveform [--debug] [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
	println()
}
