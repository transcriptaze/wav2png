package wav2png

import (
	"bytes"
	_ "embed"
	"image/png"
	"reflect"
	"testing"
)

//go:embed ice.png
var file []byte

func TestPaletteFromPng(t *testing.T) {
	buffer := bytes.NewBuffer(file)
	img, err := png.Decode(buffer)
	if err != nil {
		t.Fatalf("Error loading PNG file (%v)", err)
	}

	palette, err := PaletteFromPng(img)
	if err != nil {
		t.Fatalf("Error creating palette from PNG (%v)", err)
	} else if palette == nil {
		t.Fatalf("Failed to create palette from PNG (%v)", palette)
	}

	if !reflect.DeepEqual(*palette, ice) {
		if len(palette.colours) != len(ice.colours) {
			t.Errorf("Palette size is incorrect - expected:%v, got:%v", len(ice.colours), len(palette.colours))
		} else {
			for i, colour := range ice.colours {
				if !reflect.DeepEqual(palette.colours[i], colour) {
					t.Errorf("Colour %v is incorrect - expected:%v, got:%v", i, ice.colours[i], palette.colours[i])
				}
			}
		}
	}
}
