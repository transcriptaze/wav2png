package options

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/transcriptaze/wav2png/cmd/options"
	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/wav2png"
)

var defaults = settings{
	Size: options.Size{
		Width:  645,
		Height: 390,
	},
	Padding: options.Padding(2),
	Palette: "ice",
	Fill: options.Fill{
		Fill:   "solid",
		Colour: "#000000",
		Alpha:  255,
	},
	Grid: options.Grid{
		Grid:   "square",
		Colour: "#008000",
		Alpha:  255,
		Size:   "~64",
		WH:     "~64x48",
	},
	Antialias: options.Antialias{
		Type: "vertical",
	},
	Scale: options.Scale{
		Horizontal: 1.0,
		Vertical:   1.0,
	},
	Style: "default",
}

type Options struct {
	WAV     string
	PNG     string
	Height  uint
	Width   uint
	Padding int
	From    *time.Duration
	To      *time.Duration
	Mix     options.Mix

	Palette   wav2png.Palette
	FillSpec  fills.FillSpec
	GridSpec  grids.GridSpec
	Antialias kernels.Kernel
	VScale    float64

	Style string

	Debug bool
}

func NewOptions() Options {
	return Options{
		Width:     uint(defaults.Size.Width),
		Height:    uint(defaults.Size.Height),
		Padding:   int(defaults.Padding),
		Palette:   defaults.Palette.Palette(),
		FillSpec:  defaults.Fill.FillSpec(),
		GridSpec:  defaults.Grid.GridSpec(),
		Antialias: defaults.Antialias.Kernel(),
		VScale:    defaults.Scale.Vertical,
		Style:     defaults.Style,
		Debug:     false,
	}
}

func (o *Options) Parse() error {
	// ... initialise options from .settings
	if err := defaults.Load(".settings"); err == nil {
		o.Width = uint(defaults.Size.Width)
		o.Height = uint(defaults.Size.Height)
		o.Padding = int(defaults.Padding)
		o.Palette = defaults.Palette.Palette()
		o.FillSpec = defaults.Fill.FillSpec()
		o.GridSpec = defaults.Grid.GridSpec()
		o.Antialias = defaults.Antialias.Kernel()
		o.VScale = defaults.Scale.Vertical
		o.Style = defaults.Style
	}

	// ... override default settings with command line options
	var out string
	var settings string
	var width uint
	var height uint
	var padding int
	var start time.Duration
	var end time.Duration
	var mix options.Mix
	var style string

	palette := defaults.Palette
	grid := defaults.Grid
	fill := defaults.Fill
	antialias := defaults.Antialias
	scale := defaults.Scale

	flag.StringVar(&settings, "settings", ".settings", "JSON file with the default settings")
	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.UintVar(&width, "width", o.Width, "Image width (pixels)")
	flag.UintVar(&height, "height", o.Height, "Image height (pixels)")
	flag.IntVar(&padding, "padding", o.Padding, "Image padding (pixels)")
	flag.Var(&palette, "palette", "name of built-in palette or PNG file")
	flag.Var(&grid, "grid", "'grid' specification")
	flag.Var(&fill, "fill", "'fill' specification")
	flag.Var(&antialias, "antialias", "'antialias' specification")
	flag.Var(&scale, "scale", "vertical scaling")
	flag.DurationVar(&start, "start", 0, "start time of audio selection")
	flag.DurationVar(&end, "end", 1*time.Hour, "end time of audio selection")
	flag.Var(&mix, "mix", "channel mix")
	flag.StringVar(&style, "style", "", "render style")
	flag.BoolVar(&o.Debug, "debug", false, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return fmt.Errorf("missing WAV file")
	}

	wavfile := path.Clean(flag.Arg(0))
	filename := path.Base(wavfile)
	ext := path.Ext(filename)
	png := strings.TrimSuffix(filename, ext) + ".png"

	initialise := func(a *flag.Flag) {
		switch a.Name {
		case "out":
			info, err := os.Stat(out)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("\n   ERROR: %v\n\n", err)
				os.Exit(1)
			} else if err == nil && info.IsDir() {
				png = path.Join(out, png)
			} else {
				png = out
			}

		case "settings":
			if err := defaults.Load(settings); err == nil {
				o.Width = uint(defaults.Size.Width)
				o.Height = uint(defaults.Size.Height)
				o.Padding = int(defaults.Padding)
				o.Palette = defaults.Palette.Palette()
				o.FillSpec = defaults.Fill.FillSpec()
				o.GridSpec = defaults.Grid.GridSpec()
				o.Antialias = defaults.Antialias.Kernel()
				o.VScale = defaults.Scale.Vertical
			}

		case "start":
			o.From = &start

		case "end":
			o.To = &end

		case "mix":
			o.Mix = mix
		}
	}

	overrides := func(a *flag.Flag) {
		switch a.Name {
		case "width":
			o.Width = width

		case "height":
			o.Height = height

		case "padding":
			o.Padding = padding

		case "palette":
			o.Palette = palette.Palette()

		case "fill":
			o.FillSpec = fill.FillSpec()

		case "grid":
			o.GridSpec = grid.GridSpec()

		case "antialias":
			o.Antialias = antialias.Kernel()

		case "scale":
			o.VScale = scale.Vertical

		case "style":
			o.Style = style
		}
	}

	flag.Visit(initialise)
	flag.Visit(overrides)

	o.WAV = wavfile
	o.PNG = png

	return nil
}
