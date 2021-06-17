package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
	"strings"
	"time"

	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/wav2png"
)

type audio struct {
	sampleRate float64
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
	from    *time.Duration
	to      *time.Duration
}{
	palette: wav2png.Ice,
}

func main() {
	var out string
	var height uint
	var width uint
	var padding int
	var debug bool

	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.UintVar(&height, "height", 390, "Image height (pixels)")
	flag.UintVar(&width, "width", 645, "Image width (pixels)")
	flag.IntVar(&padding, "padding", 0, "Image padding (pixels)")
	flag.BoolVar(&debug, "debug", false, "Displays diagnostic information")
	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}

	wavfile := path.Clean(flag.Arg(0))

	filename := path.Base(wavfile)
	ext := path.Ext(filename)
	pngfile := strings.TrimSuffix(filename, ext) + ".png"
	if out != "" {
		info, err := os.Stat(out)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("\n   ERROR: %v\n\n", err)
			os.Exit(1)
		} else if err == nil && info.IsDir() {
			pngfile = path.Join(out, pngfile)
		} else {
			pngfile = out
		}
	}

	w, err := read(wavfile)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	} else if w == nil {
		fmt.Printf("\n   ERROR: unable to read WAV file\n")
		os.Exit(1)
	}

	if debug {
		fmt.Println()
		fmt.Printf("   File:        %v\n", wavfile)
		// fmt.Printf("   Format:   %+v\n", *format)
		// fmt.Printf("   Bits:     %+v\n", bits)
		fmt.Printf("   Sample Rate: %v\n", w.sampleRate)
		fmt.Printf("   Duration:    %v\n", w.duration)
		fmt.Printf("   Samples:     %v\n", w.length)
		fmt.Printf("   PNG:         %v\n", pngfile)
		fmt.Println()
	}

	img, err := render(*w, settings)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	if err := write(img, pngfile); err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}

	// params := wav2png.Params{
	// 	Width:   width,
	// 	Height:  height,
	// 	Padding: padding,
	// }
	//
	// wav2png.Draw(wavfile, pngfile, params)
}

func render(wav audio, settings Settings) (*image.NRGBA, error) {
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

	// if cache.from != nil {
	// 	v := int(math.Floor(cache.from.Seconds() * fs))
	// 	if v > 0 && v <= len(wav.samples) {
	// 		start = v
	// 	}
	// }

	// if cache.to != nil {
	// 	v := int(math.Floor(cache.to.Seconds() * fs))
	// 	if v < start {
	// 		end = start
	// 	} else if v <= len(wav.samples) {
	// 		end = v
	// 	}
	// }

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
