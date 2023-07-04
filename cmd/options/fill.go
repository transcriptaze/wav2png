package options

import (
	"fmt"
	"image/color"
	"regexp"
	"strings"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Fill struct {
	Fill   string `json:"fill"`
	Colour string `json:"colour"`
	Alpha  uint8  `json:"alpha"`
}

func (f Fill) String() string {
	if f.Fill == "solid" {
		return fmt.Sprintf("%v:%v%02x", f.Fill, f.Colour, f.Alpha)
	}

	return "??"
}

func (f *Fill) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(solid).*").FindStringSubmatch(ss)

	if len(match) > 1 {
		switch match[1] {
		case "solid":
			f.Fill = "solid"

			match = regexp.MustCompile("^solid:(#[[:xdigit:]]{8}).*").FindStringSubmatch(ss)
			if len(match) > 1 {
				color := colour(match[1])
				f.Colour = fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
				f.Alpha = color.A
			}
		}
	}

	return nil
}

func (f *Fill) FillSpec() wav2png.FillSpec {
	colour := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}

	red := uint8(0)
	green := uint8(0)
	blue := uint8(0)
	alpha := f.Alpha
	if _, err := fmt.Sscanf(f.Colour, "#%02x%02x%02x", &red, &green, &blue); err == nil {
		colour = color.NRGBA{R: red, G: green, B: blue, A: alpha}
	}

	return wav2png.NewSolidFill(colour)
}
