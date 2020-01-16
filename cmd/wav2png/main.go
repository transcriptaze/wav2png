package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"wav2png/wav2png"
)

func main() {
	var out = ""
	var height = uint(256)
	var width = uint(1024)
	var padding = uint(0)

	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.UintVar(&height, "height", 256, "Image height (pixels)")
	flag.UintVar(&width, "width", 1024, "Image width (pixels)")
	flag.UintVar(&padding, "padding", 0, "Image padding (pixels)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		println()
		println("   Usage: waveform [--height <height>] [--width <width>] [--padding <padding>] [--out <filepath>] <filename>")
		println()

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

	err := wav2png.Plot(wavfile, pngfile, width, height, padding)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n", err)
		os.Exit(1)
	}
}
