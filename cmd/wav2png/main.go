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

const VERSION = "v1.1.0"

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

	if len(os.Args) > 1 && os.Args[1] == "help" {
		help()
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

	from := 0 * time.Second
	to := audio.duration

	if options.From != nil {
		from = *options.From
	}

	if options.To != nil {
		to = *options.To
	}

	if options.Debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", options.WAV)
		fmt.Printf("   Channels:    %v\n", audio.channels)
		fmt.Printf("   Format:      %v\n", audio.format)
		fmt.Printf("   Sample Rate: %v\n", audio.sampleRate)
		fmt.Printf("   Duration:    %v\n", audio.duration)
		fmt.Printf("   Samples:     %v\n", audio.length)
		fmt.Printf("   PNG:         %v\n", options.PNG)
		fmt.Println()
	}

	img, err := render(*audio, from, to, options)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	if err := write(img, options.PNG); err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}
}

func render(wav audio, from, to time.Duration, options options.Options) (*image.NRGBA, error) {
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

	x0 := padding
	y0 := padding
	x1 := x0 + w
	y1 := y0 + h

	origin := image.Pt(0, 0)
	rect := image.Rect(x0, y0, x1, y1)
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
	fmt.Println("   Usage: wav2png [--debug] [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
	fmt.Println()
}

func help() {
	fmt.Println()
	fmt.Println("   Usage: wav2png [--debug] [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
	fmt.Println()
	fmt.Println()
	fmt.Println("       <wav>         WAV file to render.")
	fmt.Println()
	fmt.Println("       --out <path>  File path for MP4 file - if <path> is a directory, the WAV file name is")
	fmt.Println("                     used and defaults to the WAV file base path. wav2mp4 generates a set of ffmpeg frames ")
	fmt.Println("                     files in the 'frames' subdirectory of the out file directory. ")
	fmt.Println()
	fmt.Println("       --debug       Displays occasionally useful diagnostic information.")
	fmt.Println()
	fmt.Println()
	fmt.Println("   Options:")
	fmt.Println()
	fmt.Println("    --settings <file>      JSON file with the default settings for the height, width, etc.Defaults to .settings.json if")
	fmt.Println("                           not specified, falling back to internal default values if .settings.json does not exist")
	fmt.Println()
	fmt.Println("    --width    <pixels>    Width (in pixels) of the PNG image. Valid values are in the range 32 to 8192, defaults to")
	fmt.Println("                           645px")
	fmt.Println()
	fmt.Println("    --height   <pixels>    Height (in pixels) of the PNG image. Valid values are in the range 32 to 8192, defaults to")
	fmt.Println("                           395px")
	fmt.Println()
	fmt.Println("    --padding  <pixels>    Padding (in pixels) between the border of the PNG and the extent of the rendered waveform. Valid")
	fmt.Println("                           values are in the range -16 to 32, defaults to 2")
	fmt.Println()
	fmt.Println("    --palette  <palette>   Palette used to colour the waveform. May be the name of one of the internal colour palettes")
	fmt.Println("                           or a user provided PNG file. Defaults to 'ice'")
	fmt.Println()
	fmt.Println("    --fill <fillspec>      Fill specification for the background colour, in the form type:colour e.g. solid:#0000ffff")
	fmt.Println()
	fmt.Println("    --grid <gridspec>      Grid specification for an optional rectilinear grid, in the form type:colour:size:overlay, e.g.")
	fmt.Println("                           - none")
	fmt.Println("                           - square:#008000ff:~64")
	fmt.Println("                           - rectangle:#008000ff:~64x48:overlay")
	fmt.Println()
	fmt.Println("                           The size may preceded by a 'fit':")
	fmt.Println("                           - ~  approximate")
	fmt.Println("                           - =  exact")
	fmt.Println("                           - ≥  at least")
	fmt.Println("                           - ≤  at most")
	fmt.Println("                           - >  greater than")
	fmt.Println("                           - <  less than")
	fmt.Println()
	fmt.Println("                           If gridspec includes :overlay, the grid is rendered 'in front' of the waveform.")
	fmt.Println()
	fmt.Println("                           The default gridspec is 'square:#008000ff:~64'")
	fmt.Println()
	fmt.Println("    --antialias <kernel>   The antialising kernel applied to soften the rendered PNG. Valid values are:")
	fmt.Println("                           - none        no antialiasing")
	fmt.Println("                           - horizontal  blurs horizontal edges")
	fmt.Println("                           - vertical    blurs vertical edges")
	fmt.Println("                           - soft        blurs both horizontal and vertical edges")
	fmt.Println()
	fmt.Println("                           The default kernel is 'vertical'")
	fmt.Println()
	fmt.Println("    --scale <scale>        A vertical scaling factor to size the height of the rendered waveform. The valid range")
	fmt.Println("                           is 0.2 to 5.0, defaults to 1.0.")
	fmt.Println()
	fmt.Println("    --mix  <mixspec>       Specifies how to combine channels from a stereo WAV file. Valid values are:")
	fmt.Println("                           - 'L'    Renders the left channel only")
	fmt.Println("                           - 'R'    Renders the right channel only")
	fmt.Println("                           - 'L+R'  Combines the left and right channels")
	fmt.Println()
	fmt.Println("                           Defaults to 'L+R'.")
	fmt.Println()
	fmt.Println("    --start <time>         The start time of the segment of audio to render, in Go time format (e.g. 10s or 1m5s).")
	fmt.Println("                           Defaults to 0s.")
	fmt.Println()
	fmt.Println("    --end <time>           The end time of the segment of audio to render, in Go time format (e.g. 10s or 1m5s)")
	fmt.Println("                           Defaults to the end of the audio.")
	fmt.Println()
}

func version() {
	fmt.Println()
	fmt.Printf("   wav2png  %v\n", VERSION)
	fmt.Println()
}
