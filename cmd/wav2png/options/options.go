package options

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/transcriptaze/wav2png/wav2png"
)

var defaults = settings{
	Size: Size{
		Width:  645,
		Height: 390,
	},
	Padding: Padding(2),
	Palette: "ice",
	Fill: Fill{
		Fill:   "solid",
		Colour: "#000000",
		Alpha:  255,
	},
	Grid: Grid{
		Grid:   "square",
		Colour: "#008000",
		Alpha:  255,
		Size:   "~64",
		WH:     "~64x48",
	},
	Antialias: Antialias{
		Type: "vertical",
	},
	Scale: Scale{
		Horizontal: 1.0,
		Vertical:   1.0,
	},
}

type Options struct {
	WAV     string
	PNG     string
	Height  uint
	Width   uint
	Padding int
	From    *time.Duration
	To      *time.Duration
	Mix     Mix

	Palette   wav2png.Palette
	FillSpec  wav2png.FillSpec
	GridSpec  wav2png.GridSpec
	Antialias wav2png.Kernel
	VScale    float64

	Debug bool
}

func NewOptions() Options {
	return Options{
		Width:     uint(defaults.Size.Width),
		Height:    uint(defaults.Size.Height),
		Padding:   int(defaults.Padding),
		Palette:   defaults.Palette.palette(),
		FillSpec:  defaults.Fill.fillspec(),
		GridSpec:  defaults.Grid.gridspec(),
		Antialias: defaults.Antialias.kernel(),
		VScale:    defaults.Scale.Vertical,
		Debug:     false,
	}
}

func (o *Options) Parse() error {
	// ... initialise options from .settings
	if err := defaults.Load(".settings"); err == nil {
		o.Width = uint(defaults.Size.Width)
		o.Height = uint(defaults.Size.Height)
		o.Padding = int(defaults.Padding)
		o.Palette = defaults.Palette.palette()
		o.FillSpec = defaults.Fill.fillspec()
		o.GridSpec = defaults.Grid.gridspec()
		o.Antialias = defaults.Antialias.kernel()
		o.VScale = defaults.Scale.Vertical
	}

	// ... override default settings with command line options
	var out string
	var settings string
	var width uint
	var height uint
	var padding int
	var start time.Duration
	var end time.Duration
	var mix Mix

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
	flag.BoolVar(&o.Debug, "debug", false, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return fmt.Errorf("Missing WAV file")
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
				o.Palette = defaults.Palette.palette()
				o.FillSpec = defaults.Fill.fillspec()
				o.GridSpec = defaults.Grid.gridspec()
				o.Antialias = defaults.Antialias.kernel()
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
			o.Palette = palette.palette()

		case "fill":
			o.FillSpec = fill.fillspec()

		case "grid":
			o.GridSpec = grid.gridspec()

		case "antialias":
			o.Antialias = antialias.kernel()

		case "scale":
			o.VScale = scale.Vertical
		}
	}

	flag.Visit(initialise)
	flag.Visit(overrides)

	o.WAV = wavfile
	o.PNG = png

	return nil
}
