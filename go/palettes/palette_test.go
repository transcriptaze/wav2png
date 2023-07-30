package palettes

import (
	"bytes"
	_ "embed"
	"image/png"
	"reflect"
	"testing"
)

//go:embed palettes/ice.png
var ice []byte

func TestPaletteFromPng(t *testing.T) {
	buffer := bytes.NewBuffer(ice)
	img, err := png.Decode(buffer)
	if err != nil {
		t.Fatalf("Error loading PNG file (%v)", err)
	}

	palette, err := PaletteFromPng("test", img)
	if err != nil {
		t.Fatalf("Error creating palette from PNG (%v)", err)
	} else if palette == nil {
		t.Fatalf("Failed to create palette from PNG (%v)", palette)
	}

	if !reflect.DeepEqual(*palette, Ice) {
		if len(palette.colours) != len(Ice.colours) {
			t.Errorf("Palette size is incorrect - expected:%v, got:%v", len(Ice.colours), len(palette.colours))
		} else {
			for i, colour := range Ice.colours {
				if !reflect.DeepEqual(palette.colours[i], colour) {
					t.Errorf("Colour %v is incorrect - expected:%v, got:%v", i, Ice.colours[i], palette.colours[i])
				}
			}
		}
	}
}
