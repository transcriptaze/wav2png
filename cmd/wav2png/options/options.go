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
	Debug   bool
	From    *time.Duration
	To      *time.Duration

	FillSpec  wav2png.FillSpec
	GridSpec  wav2png.GridSpec
	Antialias wav2png.Kernel
	VScale    float64
}

func NewOptions() Options {
	return Options{
		Width:   uint(defaults.Size.width),
		Height:  uint(defaults.Size.height),
		Padding: int(defaults.Padding),

		FillSpec:  defaults.Fill.fillspec(),
		GridSpec:  defaults.Grid.gridspec(),
		Antialias: defaults.Antialias.kernel(),
		VScale:    defaults.Scale.Vertical,

		Debug: false,
	}
}

func (o *Options) Parse() error {
	var out string
	var start time.Duration
	var end time.Duration
	grid := defaults.Grid
	fill := defaults.Fill
	antialias := defaults.Antialias

	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.UintVar(&o.Height, "height", 390, "Image height (pixels)")
	flag.UintVar(&o.Width, "width", 645, "Image width (pixels)")
	flag.IntVar(&o.Padding, "padding", 0, "Image padding (pixels)")
	flag.Var(&grid, "grid", "'grid' specification")
	flag.Var(&fill, "fill", "'fill' specification")
	flag.Var(&antialias, "antialias", "'antialias' specification")
	flag.DurationVar(&start, "start", 0, "start time of audio selection")
	flag.DurationVar(&end, "end", 1*time.Hour, "end time of audio selection")
	flag.BoolVar(&o.Debug, "debug", false, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return fmt.Errorf("Missing WAV file")
	}

	wavfile := path.Clean(flag.Arg(0))
	filename := path.Base(wavfile)
	ext := path.Ext(filename)
	png := strings.TrimSuffix(filename, ext) + ".png"

	visitor := func(a *flag.Flag) {
		switch a.Name {
		case "fill":
			o.FillSpec = fill.fillspec()

		case "grid":
			o.GridSpec = grid.gridspec()

		case "antialias":
			o.Antialias = antialias.kernel()

		case "start":
			o.From = &start

		case "end":
			o.To = &end

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
		}
	}

	flag.Visit(visitor)

	o.WAV = wavfile
	o.PNG = png

	return nil
}
