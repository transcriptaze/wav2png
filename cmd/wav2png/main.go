package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/transcriptaze/wav2png/cmd/options"
	"github.com/transcriptaze/wav2png/compositor"
	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/styles"
)

const VERSION = "v1.1.0"

var opts = struct {
	start time.Duration
	end   time.Duration
	mix   options.Mix

	out     string
	width   uint
	height  uint
	padding int

	debug bool
}{
	out:    "",
	width:  800,
	height: 600,
	debug:  false,
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
	var style styles.Style
	var err error

	exit := func(err error) {
		fmt.Printf("   *** ERROR: %v", err)
		usage()
		os.Exit(1)
	}

	if wavfile, err = parse(); err != nil {
		exit(err)
	} else if outfile, err = makeOutFile(wavfile); err != nil {
		exit(err)
	} else if style, err = makeStyle(); err != nil {
		exit(err)
	} else if audio, err = getAudio(wavfile); err != nil {
		exit(err)
	}

	if img, err := render(audio, style); err != nil {
		exit(err)
	} else if err := write(img, outfile); err != nil {
		exit(err)
	}
}

func parse() (string, error) {
	flag.StringVar(&opts.out, "out", opts.out, "Output file (or directory)")
	flag.UintVar(&opts.width, "width", opts.width, "Image width (pixels)")
	flag.UintVar(&opts.height, "height", opts.height, "Image height (pixels)")
	flag.IntVar(&opts.padding, "padding", opts.padding, "Image padding (pixels)")
	// flag.Var(&palette, "palette", "name of built-in palette or PNG file")
	// flag.Var(&grid, "grid", "'grid' specification")
	// flag.Var(&fill, "fill", "'fill' specification")
	// flag.Var(&antialias, "antialias", "'antialias' specification")
	// flag.Var(&scale, "scale", "vertical scaling")
	flag.DurationVar(&opts.start, "start", 0, "start time of audio selection")
	flag.DurationVar(&opts.end, "end", 1*time.Hour, "end time of audio selection")
	flag.Var(&opts.mix, "mix", "channel mix")
	// flag.StringVar(&style, "style", "", "render style")
	flag.BoolVar(&opts.debug, "debug", opts.debug, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return "", fmt.Errorf("missing WAV file")
	} else {
		return filepath.Clean(flag.Args()[0]), nil
	}
}

func makeOutFile(wavfile string) (png string, err error) {
	filename := filepath.Base(wavfile)
	ext := filepath.Ext(filename)
	png = strings.TrimSuffix(filename, ext) + ".png"

	var info fs.FileInfo

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "out" {
			info, err = os.Stat(opts.out)
			if err != nil && !os.IsNotExist(err) {
				return
			} else if err == nil && info.IsDir() {
				png = filepath.Join(opts.out, png)
			} else {
				png = opts.out
			}
		}
	})

	return
}

func makeStyle() (style styles.Style, err error) {
	style = styles.NewStyle()
	err = nil

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "width":
			style.Width = opts.width
		case "height":
			style.Height = opts.height
		case "padding":
			style.Padding = opts.padding
		}
	})

	return
}

func getAudio(file string) (pcm []float32, err error) {
	var f *os.File
	var audio encoding.Audio

	if f, err = os.Open(file); err != nil {
		return
	}

	defer f.Close()

	if audio, err = encoding.Decode(f); err != nil {
		return
	}

	from := 0 * time.Second
	to := audio.Duration

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "start" && opts.start < audio.Duration {
			from = opts.start
		} else if f.Name == "end" && opts.end < audio.Duration {
			to = opts.end
		}
	})

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

	fs := audio.SampleRate
	samples := mix(audio, opts.mix.Channels()...)
	start := int(math.Floor(from.Seconds() * fs))
	end := int(math.Floor(to.Seconds() * fs))

	pcm = samples[start:end]

	return
}

func render(audio []float32, style styles.Style) (*image.NRGBA, error) {
	compositor := compositor.FromStyle(style)

	return compositor.Render(audio)
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
