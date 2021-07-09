package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/transcriptaze/wav2png/cmd/wav2mp4/options"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/wav2png"
)

const VERSION = "v1.0.0"

type audio struct {
	sampleRate float64
	format     string
	channels   int
	duration   time.Duration
	length     int
	samples    [][]float32
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		version()
		os.Exit(0)
	}

	options := options.NewOptions()
	if err := options.Parse(); err != nil {
		usage()
		os.Exit(1)
	}

	audio, err := read(options.WAV)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	} else if audio == nil {
		fmt.Printf("\n   ERROR: unable to read WAV file\n")
		os.Exit(1)
	}

	if options.Debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", options.WAV)
		fmt.Printf("   Channels:    %v\n", audio.channels)
		fmt.Printf("   Format:      %v\n", audio.format)
		fmt.Printf("   Sample Rate: %v\n", audio.sampleRate)
		fmt.Printf("   Duration:    %v\n", audio.duration)
		fmt.Printf("   Samples:     %v\n", audio.length)
		fmt.Printf("   MP4:         %v\n", options.MP4)
		fmt.Printf("   - window:    %v\n", options.Window)
		fmt.Printf("   - FPS:       %v\n", options.FPS)
		fmt.Println()
	}

	from := 0 * time.Second
	to := audio.duration

	if options.From != nil {
		from = *options.From
	}

	if options.To != nil {
		to = *options.To
	}

	if err := os.MkdirAll(options.Frames, 0777); err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	w := int(options.Width)
	h := int(options.Height)
	padding := options.Padding
	duration := to - from

	x0 := 0
	y0 := 0
	if padding > 0 {
		x0 = padding
		y0 = padding
		w -= 2 * padding
		h -= 2 * padding
	}

	cursor := options.Cursor.Render(h)
	fn := options.Cursor.Fn()
	frames := int(math.Floor((to - from).Seconds() * options.FPS))

	for frame := 0; frame <= frames; frame++ {
		png := filepath.Join(options.Frames, fmt.Sprintf("frame-%05v.png", frame+1))

		percentage := float64(frame) / float64(frames)
		t := percentage * duration.Seconds()
		offset, x, shift := fn(seconds(t), options.Window, duration)

		//		offset := seconds((to - from - options.Window).Seconds() * float64(frame) / float64(frames))

		start := from + offset
		end := start + options.Window
		if end > audio.duration {
			end = audio.duration
			start = end - options.Window
		}

		img, err := render(*audio, start, end, options, shift)
		if err != nil {
			fmt.Printf("\n   ERROR: %v\n", err)
			os.Exit(1)
		} else if img == nil {
			fmt.Printf("\n   ERROR: error creating frame\n")
			os.Exit(1)
		}

		if cursor != nil {
			cx := x0 + int(math.Round(x*float64(w-1)))
			cy := y0
			cw := cursor.Bounds().Dx()
			ch := cursor.Bounds().Dy()
			cxy := image.Pt(cw/2, 0)

			draw.Draw(img, image.Rect(cx, cy, cx+cw, cy+ch), cursor, cxy, draw.Over)
		}

		if options.Debug {
			fmt.Printf("   ... frame %-5d  %-8v %-8v %v\n", frame+1, start.Round(time.Millisecond), end.Round(time.Millisecond), png)
		}

		if err := write(img, png); err != nil {
			fmt.Printf("\n   ERROR: %v\n", err)
			os.Exit(1)
		}
	}
}

func render(wav audio, from, to time.Duration, options options.Options, shift float64) (*image.NRGBA, error) {
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
		w = width - 2*padding
		h = height - 2*padding
	}

	fs := wav.sampleRate
	samples := mix(wav, options.Mix.Channels()...)

	start := int(math.Floor(from.Seconds() * fs))
	if start < 0 || start > len(samples) {
		return nil, fmt.Errorf("Start position not in range %v-%v", from, wav.duration)
	}

	end := int(math.Floor(to.Seconds() * fs))
	if end < 0 || end < start || end > len(samples) {
		return nil, fmt.Errorf("End position not in range %v-%v", from, wav.duration)
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := wav2png.Grid(gridspec, width, height, padding)
	waveform := wav2png.Render(samples[start:end], fs, w, h, palette, vscale)
	antialiased := wav2png.Antialias(waveform, kernel)

	offset := int(math.Round(shift * float64(w)))
	x0 := padding + offset
	x1 := padding + w
	if x0 < padding {
		x0 = padding
		x1 = padding + offset + w
	}

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, padding, x1, h)
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
	fmt.Println()
	fmt.Println("   Usage: waveform [--debug] [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
	fmt.Println()
}

func version() {
	fmt.Println()
	fmt.Printf("   wav2png  %v\n", VERSION)
	fmt.Println()
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
