package main

import (
	"flag"
	"fmt"
	"github.com/go-audio/wav"
	"os"
	"path"
	"strings"
	"wav2png/wav2png"
)

func main() {
	var out = ""

	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		println()
		println("   Usage: waveform [--out <filepath>] <filename>")
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

	file, err := os.Open(wavfile)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	defer file.Close()

	decoder := wav.NewDecoder(file)

	decoder.FwdToPCM()

	format := decoder.Format()
	bits := decoder.SampleBitDepth()
	duration, err := decoder.Duration()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Printf("   File:     %s\n", wavfile)
	fmt.Printf("   Format:   %+v\n", *format)
	fmt.Printf("   Bits:     %+v\n", bits)
	fmt.Printf("   Duration: %s\n", duration)
	fmt.Printf("   Length:   %d\n", decoder.PCMLen())

	err = wav2png.Plot(decoder, pngfile)
	if err != nil {
		fmt.Printf("\n   ERROR: EOF\n\n")
		os.Exit(1)
	}
}
