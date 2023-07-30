package palettes

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	// "os"
	// "regexp"
	// "strings"
)

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

type Palette struct {
	name    string
	colours []color.NRGBA
}

func PaletteFromPng(name string, png image.Image) (*Palette, error) {
	bounds := png.Bounds()
	if bounds.Empty() {
		return nil, fmt.Errorf("cannot create palette from empty PNG")
	}

	h := bounds.Size().Y
	colours := make([]color.NRGBA, h)
	nrgba := color.NRGBAModel

	for i := 0; i < h; i++ {
		colours[i] = nrgba.Convert(png.At(0, i)).(color.NRGBA)
	}

	return &Palette{
		name:    name,
		colours: colours,
	}, nil
}

func (p Palette) String() string {
	return p.name
}

// func (p *Palette) Set(s string) error {
// ss := strings.ToLower(s)
// match := regexp.MustCompile("^(ice|fire|aurora|horizon|amber|blue|green|gold)$").FindStringSubmatch(ss)

// if len(match) > 1 {
// 	*p = Palette(match[1])
// 	return nil
// }

// if info, err := os.Stat(s); os.IsNotExist(err) {
// 	return fmt.Errorf("Palette %v does not exist", s)
// } else if info.Mode().IsDir() || !info.Mode().IsRegular() {
// 	return fmt.Errorf("Palette file %v is not a file", s)
// } else {
// 	*p = Palette(s)
// }

// return nil
// }

func (p *Palette) Realize() []color.NRGBA {
	return p.colours
}

// func (p Palette) Palette() Palette {
// if b, ok := palettes[string(p)]; ok {
// 	if v := p.decode(b); v != nil {
// 		return *v
// 	}
// }
//
// if b, err := os.ReadFile(string(p)); err == nil {
// 	if v := p.decode(b); v != nil {
// 		return *v
// 	}
// }
//
// return Ice
// }

func decode(name string, b []byte) *Palette {
	if img, err := png.Decode(bytes.NewBuffer(b)); err != nil {
		return nil
	} else if palette, err := PaletteFromPng(name, img); err != nil {
		return nil
	} else {
		return palette
	}
}

func mustDecode(name string, b []byte) Palette {
	if img, err := png.Decode(bytes.NewBuffer(b)); err != nil {
		panic(fmt.Errorf("invalid palette (%v)", err))
	} else if palette, err := PaletteFromPng(name, img); err != nil {
		panic(fmt.Errorf("invalid palette (%v)", err))
	} else {
		return *palette
	}
}

var Ice = mustDecode("ice", palette_ice)
var Fire = mustDecode("fire", palette_fire)
var Aurora = mustDecode("aurora", palette_aurora)
var Horizon = mustDecode("horizon", palette_horizon)
var Amber = mustDecode("amber", palette_amber)
var Blue = mustDecode("blue", palette_blue)
var Green = mustDecode("green", palette_green)
var Gold = mustDecode("gold", palette_gold)
