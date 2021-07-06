package options

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Palette string

//go:embed fire.png
var fire []byte

//go:embed aurora.png
var aurora []byte

//go:embed horizon.png
var horizon []byte

//go:embed amber.png
var amber []byte

//go:embed blue.png
var blue []byte

//go:embed green.png
var green []byte

//go:embed gold.png
var gold []byte

var palettes = map[string][]byte{
	"fire":    fire,
	"aurora":  aurora,
	"horizon": horizon,
	"amber":   amber,
	"blue":    blue,
	"green":   green,
	"gold":    gold,
}

func (p Palette) String() string {
	return fmt.Sprintf("%v", string(p))
}

func (p *Palette) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(ice|fire|aurora|horizon|amber|blue|green|gold)$").FindStringSubmatch(ss)

	if match != nil && len(match) > 1 {
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
