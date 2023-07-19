package palettes

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/transcriptaze/wav2png/go/wav2png"
)

type Palette string

//go:embed palettes/ice.png
var palette_ice []byte

//go:embed palettes/fire.png
var palette_fire []byte

//go:embed palettes/aurora.png
var palette_aurora []byte

//go:embed palettes/horizon.png
var palette_horizon []byte

//go:embed palettes/amber.png
var palette_amber []byte

//go:embed palettes/blue.png
var palette_blue []byte

//go:embed palettes/green.png
var palette_green []byte

//go:embed palettes/gold.png
var palette_gold []byte

var palettes = map[string][]byte{
	"ice":     palette_ice,
	"fire":    palette_fire,
	"aurora":  palette_aurora,
	"horizon": palette_horizon,
	"amber":   palette_amber,
	"blue":    palette_blue,
	"green":   palette_green,
	"gold":    palette_gold,
}

var Ice = Palette("ice")
var Fire = Palette("fire")
var Aurora = Palette("aurora")
var Horizon = Palette("horizon")
var Amber = Palette("amber")
var Blue = Palette("blue")
var Green = Palette("green")
var Gold = Palette("gold")

func (p Palette) String() string {
	return fmt.Sprintf("%v", string(p))
}

func (p *Palette) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(ice|fire|aurora|horizon|amber|blue|green|gold)$").FindStringSubmatch(ss)

	if len(match) > 1 {
		*p = Palette(match[1])
		return nil
	}

	if info, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("Palette %v does not exist", s)
	} else if info.Mode().IsDir() || !info.Mode().IsRegular() {
		return fmt.Errorf("Palette file %v is not a file", s)
	} else {
		*p = Palette(s)
	}

	return nil
}

func (p Palette) Palette() wav2png.Palette {
	if b, ok := palettes[string(p)]; ok {
		if v := p.decode(b); v != nil {
			return *v
		}
	}

	if b, err := os.ReadFile(string(p)); err == nil {
		if v := p.decode(b); v != nil {
			return *v
		}
	}

	return wav2png.Ice
}

func (p Palette) decode(b []byte) *wav2png.Palette {
	img, err := png.Decode(bytes.NewBuffer(b))
	if err != nil {
		return nil
	}

	palette, err := wav2png.PaletteFromPng(img)
	if err != nil {
		return nil
	}

	return palette
}
