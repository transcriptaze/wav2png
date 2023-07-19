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

	"github.com/transcriptaze/wav2png/go/cmd/wav2mp4/options"
	"github.com/transcriptaze/wav2png/go/encoding/wav"
	"github.com/transcriptaze/wav2png/go/fills"
	"github.com/transcriptaze/wav2png/go/grids"
	"github.com/transcriptaze/wav2png/go/wav2png"
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
		fmt.Printf("\n   ERROR: %v\n", err)
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
	window := 1 * time.Second
	fps := 30.0

	if options.Window != nil {
		window = *options.Window
	}

	if options.FPS != nil {
		fps = *options.FPS
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
		fmt.Printf("   - window:    %v\n", window)
		fmt.Printf("   - FPS:       %v\n", fps)
		fmt.Println()
	}

	fn := options.Cursor.Fn()
	frames := int(math.Floor((to - from).Seconds() * fps))

	for frame := 0; frame <= frames; frame++ {
		png := filepath.Join(options.Frames, fmt.Sprintf("frame-%05v.png", frame))

		percentage := float64(frame) / float64(frames)
		t := seconds(percentage * duration.Seconds())
		x := fn(t, duration)

		p := t - seconds(x*window.Seconds())
		q := p + window
		shift := 0.0

		if p < 0 {
			shift = (0 - p).Seconds() / window.Seconds()
			p = 0 * time.Second
			q = p + window
		}

		if q > duration {
			shift = (duration - (p + window)).Seconds() / window.Seconds()
			p = duration - window
			// q = p + window
		}

		start := from + p
		end := start + window

		if start < from {
			fmt.Printf("\n   ERROR: frame %d - invalid frame 'start' (%v)\n", frame, start)
			os.Exit(1)
		}

		if end > to {
			fmt.Printf("\n   ERROR: frame %d - invalid frame 'end' (%v)\n", frame, end)
			os.Exit(1)
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
			cw := cursor.Bounds().Dx()
			ch := cursor.Bounds().Dy()
			cx := x0 + int(math.Round(x*float64(w-1))) - cw/2
			cy := y0

			draw.Draw(img, image.Rect(cx, cy, cx+cw, cy+ch), cursor, image.Pt(0, 0), draw.Over)
		}

		if options.Debug {
			fmt.Printf("   ... frame %-5d  %-8v %-8v %v   %+.3f\n", frame, start.Round(time.Millisecond), end.Round(time.Millisecond), png, shift)
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
		return nil, fmt.Errorf("start position not in range %v-%v", from, wav.duration)
	}

	end := int(math.Floor(to.Seconds() * fs))
	if end < 0 || end < start || end > len(samples) {
		return nil, fmt.Errorf("end position not in range %v-%v", from, wav.duration)
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	grid := grids.Grid(gridspec, width, height, padding)
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

	fills.Fill(img, fillspec)

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
	fmt.Println("   Usage: wav2mp4 [--debug] [options] [--out <filepath>] --window <window> --fps <frame rate> --cursor <cursorspec> <filename>")
	fmt.Println()
}

func help() {
	fmt.Println()
	fmt.Println("   Usage: wav2mp4 [--debug] [options] [--out <filepath>] --window <window> --fps <frame rate> --cursor <cursorspec> <filename>")
	fmt.Println()
	fmt.Println()
	fmt.Println("       <wav>                  WAV file to render.")
	fmt.Println()
	fmt.Println("       --out <path>           File path for MP4 file - if <path> is a directory, the WAV file name is")
	fmt.Println("                              used and defaults to the WAV file base path. wav2mp4 generates a set of ffmpeg frames ")
	fmt.Println("                              files in the 'frames' subdirectory of the out file directory. ")
	fmt.Println()
	fmt.Println("       --window <duration>    The interval allotted to a single frame, in Go time format e.g. --window 1s.")
	fmt.Println("                              The window interval must be less than the duration of the MP4.")
	fmt.Println("")
	fmt.Println("       --fps <frame rate>     Frame rate for the MP4 in frames per second e.g. --fps 30")
	fmt.Println()
	fmt.Println("       --cursor <cursorspec>  Cursor to indicate for the current play position. A cursor is specified by the image source")
	fmt.Println("                              and dynamic:")
	fmt.Println()
	fmt.Println("                              --cursor <image>:<dynamic>")
	fmt.Println()
	fmt.Println("                              where image may be:")
	fmt.Println("                              - none")
	fmt.Println("                              - red")
	fmt.Println("                              - blue")
	fmt.Println("                              - a PNG file with a custom cursor image")
	fmt.Println()
	fmt.Println("                              The cursor 'dynamic' defaults to 'sweep' if not specified, but may be one of the following:")
	fmt.Println("                              - sweep  Moves linearly from left to right over the duration of the MP4")
	fmt.Println("                              - left   Fixed on left side")
	fmt.Println("                              - right  Fixed on right side")
	fmt.Println("                              - center Fixed in center of frame")
	fmt.Println("                              - ease   Migrates from the left to center of the frame, before moving to the right side to")
	fmt.Println("                                       finish")
	fmt.Println("                              - erf    Moves 'sigmoidally' from left to right over the duration of the MP4, with the ")
	fmt.Println("                                       sigmoid defined by the inverse error function")
	fmt.Println()
	fmt.Println("       --debug                Displays occasionally useful diagnostic information.")
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
	fmt.Printf("   wav2mp4  %v\n", VERSION)
	fmt.Println()
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
