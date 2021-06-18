package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type options struct {
	wav     string
	png     string
	height  uint
	width   uint
	padding int
	debug   bool
	from    *time.Duration
	to      *time.Duration
}

func (o *options) parse() error {
	var out string
	var start time.Duration
	var end time.Duration

	flag.StringVar(&out, "out", "", "Output file (or directory)")
	flag.UintVar(&o.height, "height", 390, "Image height (pixels)")
	flag.UintVar(&o.width, "width", 645, "Image width (pixels)")
	flag.IntVar(&o.padding, "padding", 0, "Image padding (pixels)")
	flag.DurationVar(&start, "start", 0, "start time of audio selection")
	flag.DurationVar(&end, "end", 1*time.Hour, "end time of audio selection")
	flag.BoolVar(&o.debug, "debug", false, "Displays diagnostic information")
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
		case "start":
			o.from = &start

		case "end":
			o.to = &end

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

	o.wav = wavfile
	o.png = png

	return nil
}
