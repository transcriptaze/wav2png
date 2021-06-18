package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/wav2png"
)

type audio struct {
	sampleRate float64
	format     string
	channels   int
	duration   time.Duration
	length     int
	samples    []float32
}

var settings = Settings{
	Size: Size{
		width:  645,
		height: 390,
	},

	// Palettes: Palettes{
	// 	Selected: "palette1",
	// 	Palettes: map[string][]byte{},
	// },

	Fill: Fill{
		Fill:   "solid",
		Colour: "#000000",
		Alpha:  255,
	},

	Padding: Padding(2),

	Grid: Grid{
		Grid:   "square",
		Colour: "#008000",
		Alpha:  255,
		Size:   "~64",
		WH:     "~64x48",
	},

	Antialias: Antialias{
		Type:   "vertical",
		kernel: wav2png.Vertical,
	},

	Scale: Scale{
		Horizontal: 1.0,
		Vertical:   1.0,
	},
}

var cache = struct {
	palette wav2png.Palette
}{
	palette: wav2png.Ice,
}

func main() {
	options := options{}
	if err := options.parse(); err != nil {
		usage()
		os.Exit(1)
	}

	w, err := read(options.wav)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	} else if w == nil {
		fmt.Printf("\n   ERROR: unable to read WAV file\n")
		os.Exit(1)
	}

	if options.debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", options.wav)
		fmt.Printf("   Format:      %v\n", w.format)
		fmt.Printf("   Sample Rate: %v\n", w.sampleRate)
		fmt.Printf("   Duration:    %v\n", w.duration)
		fmt.Printf("   Samples:     %v\n", w.length)
		fmt.Printf("   PNG:         %v\n", options.png)
		fmt.Println()
	}

	img, err := render(*w, options.from, options.to, settings)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	if err := write(img, options.png); err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}
}

func render(wav audio, from, to *time.Duration, settings Settings) (*image.NRGBA, error) {
	width := settings.Size.width
	height := settings.Size.height
	padding := int(settings.Padding)
	fillspec := settings.Fill.fillspec()
	gridspec := settings.Grid.gridspec()
	kernel := settings.Antialias.kernel
	vscale := settings.Scale.Vertical

	w := width
	h := height
	if padding > 0 {
		w = width - padding
		h = height - padding
	}

	start := 0
	end := len(wav.samples)
	fs := wav.sampleRate

	if from != nil {
		v := int(math.Floor(from.Seconds() * fs))
		if v > 0 && v <= len(wav.samples) {
			start = v
		}
	}

	if to != nil {
		v := int(math.Floor(to.Seconds() * fs))
		if v < start {
			end = start
		} else if v <= len(wav.samples) {
			end = v
		}
	}

	samples := wav.samples[start:end]
	duration, _ := seconds(float64(len(samples)) / fs)

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := wav2png.Grid(gridspec, width, height, padding)
	waveform := wav2png.Render(duration, samples, fs, w, h, cache.palette, vscale)
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
		length:     len(w.Samples),
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

func seconds(g float64) (time.Duration, *time.Duration) {
	t := time.Duration(g * float64(time.Second))

	return t, &t
}

func usage() {
	println()
	println("   Usage: waveform [--debug] [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
	println()
}
