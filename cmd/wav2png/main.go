package main

import (
	"fmt"
	"image"
	"image/png"
	// "math"
	"os"
	// "time"
	"flag"

	// "github.com/transcriptaze/wav2png/cmd/wav2png/options"
	"github.com/transcriptaze/wav2png/compositor"
	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/renderers/lines"
	"github.com/transcriptaze/wav2png/styles"
)

const VERSION = "v1.1.0"

var options = struct {
	out   string
	debug bool
}{
	out:   "",
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

	flag.StringVar(&options.out, "out", options.out, "Output file (or directory)")
	// flag.UintVar(&width, "width", o.Width, "Image width (pixels)")
	// flag.UintVar(&height, "height", o.Height, "Image height (pixels)")
	// flag.IntVar(&padding, "padding", o.Padding, "Image padding (pixels)")
	// flag.Var(&palette, "palette", "name of built-in palette or PNG file")
	// flag.Var(&grid, "grid", "'grid' specification")
	// flag.Var(&fill, "fill", "'fill' specification")
	// flag.Var(&antialias, "antialias", "'antialias' specification")
	// flag.Var(&scale, "scale", "vertical scaling")
	// flag.DurationVar(&start, "start", 0, "start time of audio selection")
	// flag.DurationVar(&end, "end", 1*time.Hour, "end time of audio selection")
	// flag.Var(&mix, "mix", "channel mix")
	// flag.StringVar(&style, "style", "", "render style")
	flag.BoolVar(&options.debug, "debug", options.debug, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf("   *** ERROR: missing WAV file")
		usage()
		os.Exit(1)
	}

	wavfile := flag.Args()[0]

	audio, err := read(wavfile)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	// options := options.NewOptions()
	// if err := options.Parse(); err != nil {
	// 	usage()
	// 	os.Exit(1)
	// }

	// from := 0 * time.Second
	// to := audio.Duration

	// if options.From != nil {
	// 	from = *options.From
	// }

	// if options.To != nil {
	// 	to = *options.To
	// }

	if options.debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", wavfile)
		fmt.Printf("   Channels:    %v\n", audio.Channels)
		fmt.Printf("   Format:      %v\n", audio.Format)
		fmt.Printf("   Sample Rate: %v\n", audio.SampleRate)
		fmt.Printf("   Duration:    %v\n", audio.Duration)
		fmt.Printf("   Samples:     %v\n", audio.Length)
		// fmt.Printf("   PNG:         %v\n", options.PNG)
		// fmt.Printf("   Style:       %v\n", options.Style)
		fmt.Println()
	}

	// // if _, err := styles.Load(options.Style); err != nil {
	// // 	fmt.Printf("\n   ERROR: %v\n", err)
	// // 	os.Exit(1)
	// // }

	// style := styles.LinesStyle{
	// 	Style: styles.Style{
	// 		Width:      options.Width,
	// 		Height:     options.Height,
	// 		Padding:    options.Padding,
	// 		Background: options.FillSpec,
	// 		Grid:       options.GridSpec,
	// 	},
	// 	Palette:   options.Palette,
	// 	Antialias: options.Antialias,
	// 	VScale:    options.VScale,
	// }

	// fs := audio.SampleRate
	// samples := mix(audio, options.Mix.Channels()...)
	// start := int(math.Floor(from.Seconds() * fs))
	// end := int(math.Floor(to.Seconds() * fs))

	// img, err := render(samples[start:end], style)
	// if err != nil {
	// 	fmt.Printf("\n   ERROR: %v\n", err)
	// 	os.Exit(1)
	// }

	// if err := write(img, options.PNG); err != nil {
	// 	fmt.Printf("\n   ERROR: %v\n", err)
	// 	os.Exit(1)
	// }
}

func render(audio []float32, style styles.LinesStyle) (*image.NRGBA, error) {
	compositor := compositor.NewCompositor(
		style.Width,
		style.Height,
		style.Padding,
		style.Background,
		style.Grid,
		lines.Lines{
			Palette:   style.Palette,
			AntiAlias: style.Antialias,
			VScale:    style.VScale,
		})

	return compositor.Render(audio)
}

func read(wavfile string) (encoding.Audio, error) {
	file, err := os.Open(wavfile)
	if err != nil {
		return encoding.Audio{}, err
	}

	defer file.Close()

	w, err := wav.Decode(file)
	if err != nil {
		return encoding.Audio{}, err
	}

	return encoding.Audio{
		SampleRate: float64(w.Format.SampleRate),
		Format:     fmt.Sprintf("%v", w.Format),
		Channels:   int(w.Format.Channels),
		Duration:   w.Duration(),
		Length:     w.Frames(),
		Samples:    w.Samples,
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
