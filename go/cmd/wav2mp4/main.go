package main

import (
	"flag"
	"fmt"
	"image"
	// "image/draw"
	"image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/transcriptaze/wav2png/go/cmd/options"
	"github.com/transcriptaze/wav2png/go/compositor"
	"github.com/transcriptaze/wav2png/go/encoding"
	"github.com/transcriptaze/wav2png/go/styles"
)

const VERSION = "v1.2.0"

type audio struct {
	sampleRate float64
	format     string
	channels   int
	duration   time.Duration
	length     int
	samples    [][]float32
}

var opts = struct {
	out   string
	start time.Duration
	end   time.Duration
	mix   options.Mix

	style string

	width   uint
	height  uint
	padding int

	scale styles.Scale
	fill  styles.Fill
	grid  styles.Grid

	fps    float64
	window time.Duration
	cursor options.Cursor

	debug bool
}{
	out:    "",
	style:  "",
	width:  800,
	height: 600,
	scale: styles.Scale{
		Horizontal: 1.0,
		Vertical:   1.0,
	},

	fill: styles.Fill{
		Fill:   "solid",
		Colour: "#000000",
		Alpha:  255,
	},

	grid: styles.Grid{
		Grid:   "square",
		Colour: "#008000",
		Alpha:  255,
		Size:   "~64",
		WH:     "~64x48",
	},

	fps:    30.0,
	window: 30 * time.Second,

	debug: false,
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

	var wavfile string
	var outfile string
	var audio []float32
	var fs float64
	var from time.Duration
	var to time.Duration
	var style styles.Style
	var err error

	exit := func(err error) {
		fmt.Printf("\n   *** ERROR: %v\n", err)
		usage()
		os.Exit(1)
	}

	if wavfile, err = parse(); err != nil {
		exit(err)
	} else if outfile, err = makeOutFile(wavfile); err != nil {
		exit(err)
	} else if style, err = makeStyle(); err != nil {
		exit(err)
	} else if audio, fs, from, to, err = getAudio(wavfile); err != nil {
		exit(err)
	}

	dir := filepath.Join(filepath.Dir(outfile), "frames")
	if err := os.MkdirAll(dir, 0777); err != nil {
		exit(err)
	}

	// cursor := options.Cursor.Render(h)
	fn := opts.cursor.Fn()
	duration := to - from
	window := opts.window
	fps := opts.fps
	frames := int(math.Floor(duration.Seconds() * fps))

	if opts.debug {
		fmt.Printf("   MP4:         %v\n", outfile)
		fmt.Printf("     window:    %v\n", window)
		fmt.Printf("     FPS:       %v\n", fps)
		fmt.Println()
	}

	for frame := 0; frame <= frames; frame++ {
		png := filepath.Join(dir, fmt.Sprintf("frame-%05v.png", frame))

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
		}

		start := from + p
		end := start + window

		if start < from {
			exit(fmt.Errorf("frame %d - invalid frame 'start' (%v)", frame, start))
		}

		if end > to {
			exit(fmt.Errorf("frame %d - invalid frame 'end' (%v)", frame, end))
		}

		img, err := render(audio, fs, start, end, shift, style)
		if err != nil {
			exit(err)
		} else if img == nil {
			exit(fmt.Errorf("error creating frame"))
		}

		// 	if cursor != nil {
		// 		cw := cursor.Bounds().Dx()
		// 		ch := cursor.Bounds().Dy()
		// 		cx := x0 + int(math.Round(x*float64(w-1))) - cw/2
		// 		cy := y0

		// 		draw.Draw(img, image.Rect(cx, cy, cx+cw, cy+ch), cursor, image.Pt(0, 0), draw.Over)
		// 	}

		if opts.debug {
			fmt.Printf("   ... frame %-5d  %-8v %-8v %v   %+.3f\n", frame, (from + p).Round(time.Millisecond), (from + p + window).Round(time.Millisecond), png, shift)
		}

		if err := write(img, png); err != nil {
			fmt.Printf("\n   ERROR: %v\n", err)
			os.Exit(1)
		}
	}
}

func parse() (string, error) {
	flag.StringVar(&opts.out, "out", opts.out, "Output file (or directory)")
	flag.UintVar(&opts.width, "width", opts.width, "Image width (pixels)")
	flag.UintVar(&opts.height, "height", opts.height, "Image height (pixels)")
	flag.IntVar(&opts.padding, "padding", opts.padding, "Image padding (pixels)")
	flag.Var(&opts.scale, "scale", "Vertical scaling")
	flag.StringVar(&opts.style, "style", "", "render style")
	flag.Var(&opts.fill, "fill", "(legacy) 'fill' specification")
	flag.Var(&opts.grid, "grid", "(legacy) 'grid' specification")
	flag.DurationVar(&opts.start, "start", 0, "start time of audio selection")
	flag.DurationVar(&opts.end, "end", 1*time.Hour, "end time of audio selection")
	flag.Var(&opts.mix, "mix", "channel mix")
	flag.Float64Var(&opts.fps, "fps", opts.fps, "frame rate")
	flag.DurationVar(&opts.window, "window", opts.window, "frame sample 'window'")
	flag.Var(&opts.cursor, "cursor", "name of built-in cursor or PNG file")
	flag.BoolVar(&opts.debug, "debug", opts.debug, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return "", fmt.Errorf("missing WAV file")
	} else {
		return filepath.Clean(flag.Args()[0]), nil
	}
}

func makeOutFile(wavfile string) (mp4 string, err error) {
	filename := filepath.Base(wavfile)
	ext := filepath.Ext(filename)
	mp4 = strings.TrimSuffix(filename, ext) + ".mp4"

	var info fs.FileInfo

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "out" {
			info, err = os.Stat(opts.out)
			if err != nil && !os.IsNotExist(err) {
				return
			} else if err == nil && info.IsDir() {
				mp4 = filepath.Join(opts.out, mp4)
				err = nil
			} else {
				mp4 = opts.out
				err = nil
			}
		}
	})

	return
}

func makeStyle() (style styles.Style, err error) {
	style = styles.NewStyle()
	err = nil

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "style" {
			style, err = style.Load(opts.style)
		}
	})

	if err != nil {
		return
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "width":
			style = style.WithWidth(opts.width)
		case "height":
			style = style.WithHeight(opts.height)
		case "padding":
			style = style.WithPadding(opts.padding)
		case "scale":
			style = style.WithScale(opts.scale)
		case "fill":
			style = style.WithFill(opts.fill)
		case "grid":
			style = style.WithGrid(opts.grid)
		}
	})

	return
}

func getAudio(file string) (pcm []float32, fs float64, from, to time.Duration, err error) {
	var f *os.File
	var audio encoding.Audio

	if f, err = os.Open(file); err != nil {
		return
	}

	defer f.Close()

	if audio, err = encoding.Decode(f); err != nil {
		return
	}

	if opts.debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", file)
		fmt.Printf("   Channels:    %v\n", audio.Channels)
		fmt.Printf("   Format:      %v\n", audio.Format)
		fmt.Printf("   Sample Rate: %v\n", audio.SampleRate)
		fmt.Printf("   Duration:    %v\n", audio.Duration)
		fmt.Printf("   Samples:     %v\n", audio.Length)
		fmt.Println()
	}

	pcm = mix(audio, opts.mix.Channels()...)
	fs = audio.SampleRate
	from = 0 * time.Second
	to = audio.Duration

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "start" && opts.start < audio.Duration {
			from = opts.start
		} else if f.Name == "end" && opts.end < audio.Duration {
			to = opts.end
		}
	})

	return
}

func render(audio []float32, fs float64, from, to time.Duration, shift float64, style styles.Style) (*image.NRGBA, error) {
	duration := func() time.Duration {
		return time.Duration(math.Floor(float64(len(audio))/fs)) * time.Second
	}

	start := int(math.Floor(from.Seconds() * fs))
	if start < 0 || start > len(audio) {
		return nil, fmt.Errorf("start position not in range %v-%v", from, duration())
	}

	end := int(math.Floor(to.Seconds() * fs))
	if end < 0 || end < start || end > len(audio) {
		return nil, fmt.Errorf("end position not in range %v-%v", from, duration())
	}

	compositor := compositor.FromStyle(style)

	return compositor.Render(audio[start:end])
}

func write(img *image.NRGBA, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, img)
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
