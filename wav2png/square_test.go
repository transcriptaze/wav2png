package wav2png

import (
	"image"
	"reflect"
	"testing"
)

func TestSquareGridBorder(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := image.Rect(0, 0, 640, 385)
	gridspec := SquareGrid{size: 64, padding: 0}

	border := gridspec.Border(bounds)
	if border == nil {
		t.Fatalf("GridSpec.Border unexpectedly returned %v", border)
	}

	if !reflect.DeepEqual(*border, expected) {
		t.Errorf("Incorrect border:\n   expected:%v\n   got:     %v", expected, border)
	}
}

func TestSquareGridBorderWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := image.Rect(1, 1, 641, 386)
	gridspec := SquareGrid{size: 64, padding: 1}

	border := gridspec.Border(bounds)
	if border == nil {
		t.Fatalf("GridSpec.Border unexpectedly returned %v", border)
	}

	if !reflect.DeepEqual(*border, expected) {
		t.Errorf("Incorrect border:\n   expected:%v\n   got:     %v", expected, border)
	}
}

func TestSquareGridVLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{64, 128, 192, 256, 320, 384, 448, 512, 576}
	gridspec := SquareGrid{size: 64, padding: 0}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridVLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{65, 129, 193, 257, 321, 385, 449, 513, 577}
	gridspec := SquareGrid{size: 64, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridVLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 645, 388)
	expected := []int{65, 129, 194, 258, 322, 386, 450, 515, 579}
	gridspec := SquareGrid{size: 64, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridHLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{192, 128, 64, 193, 257, 321}
	gridspec := SquareGrid{size: 64, padding: 0}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestSquareGridHLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{193, 129, 65, 194, 258, 322}
	gridspec := SquareGrid{size: 64, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestSquareGridHLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 390)
	expected := []int{194, 130, 65, 195, 260, 324}
	gridspec := SquareGrid{size: 64, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestSquareGridApproximateFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 193, 257, 321}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 383, 447, 511, 575}, hlines: []int{192, 128, 64, 192, 256, 320}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{63, 127, 190, 254, 317, 380, 444, 507, 571}, hlines: []int{189, 126, 63, 190, 253, 316}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{63, 126, 189, 252, 315, 378, 441, 504, 567}, hlines: []int{189, 126, 63, 190, 253, 316}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{65, 130, 195, 260, 325, 390, 455, 520, 585}, hlines: []int{195, 130, 65, 196, 261, 326}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 160, 96, 32, 225, 289, 353, 417}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 96, 32, 161, 225, 289}},
	}

	gridspec := SquareGrid{size: 64, padding: 0, fit: Approximate}

	for _, v := range tests {
		vlines := gridspec.VLines(v.bounds)
		hlines := gridspec.HLines(v.bounds)

		if !reflect.DeepEqual(vlines, v.vlines) {
			t.Errorf("Approximate: incorrect vertical lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.vlines, vlines)
		}

		if !reflect.DeepEqual(hlines, v.hlines) {
			t.Errorf("Approximate: incorrect horizontal lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.hlines, hlines)
		}
	}
}

func TestSquareGridExactFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 193, 257, 321}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 192, 256, 320}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 125, 61, 190, 254, 318}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 125, 61, 190, 254, 318}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{195, 131, 67, 3, 196, 260, 324, 388}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 96, 32, 161, 225, 289}},
	}

	gridspec := SquareGrid{size: 64, padding: 0, fit: Exact}

	for _, v := range tests {
		vlines := gridspec.VLines(v.bounds)
		hlines := gridspec.HLines(v.bounds)

		if !reflect.DeepEqual(vlines, v.vlines) {
			t.Errorf("Exact: incorrect vertical lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.vlines, vlines)
		}

		if !reflect.DeepEqual(hlines, v.hlines) {
			t.Errorf("Exact: incorrect horizontal lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.hlines, hlines)
		}
	}
}

func TestSquareGridAtLeastFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 193, 257, 321}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 192, 256, 320}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 125, 61, 190, 254, 318}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 125, 61, 190, 254, 318}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{65, 130, 195, 260, 325, 390, 455, 520, 585}, hlines: []int{195, 130, 65, 196, 261, 326}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 160, 96, 32, 225, 289, 353, 417}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 96, 32, 161, 225, 289}},
	}

	gridspec := SquareGrid{size: 64, padding: 0, fit: AtLeast}

	for _, v := range tests {
		vlines := gridspec.VLines(v.bounds)
		hlines := gridspec.HLines(v.bounds)

		if !reflect.DeepEqual(vlines, v.vlines) {
			t.Errorf("AtLeast: incorrect vertical lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.vlines, vlines)
		}

		if !reflect.DeepEqual(hlines, v.hlines) {
			t.Errorf("AtLeast: incorrect horizontal lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.hlines, hlines)
		}
	}
}

func TestSquareGridAtMostFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 128, 64, 193, 257, 321}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 383, 447, 511, 575}, hlines: []int{192, 128, 64, 192, 256, 320}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{63, 127, 190, 254, 317, 380, 444, 507, 571}, hlines: []int{189, 126, 63, 190, 253, 316}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{63, 126, 189, 252, 315, 378, 441, 504, 567}, hlines: []int{189, 126, 63, 190, 253, 316}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{195, 131, 67, 3, 196, 260, 324, 388}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 160, 96, 32, 225, 289, 353, 417}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 96, 32, 161, 225, 289}},
	}

	gridspec := SquareGrid{size: 64, padding: 0, fit: AtMost}

	for _, v := range tests {
		vlines := gridspec.VLines(v.bounds)
		hlines := gridspec.HLines(v.bounds)

		if !reflect.DeepEqual(vlines, v.vlines) {
			t.Errorf("AtMost: incorrect vertical lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.vlines, vlines)
		}

		if !reflect.DeepEqual(hlines, v.hlines) {
			t.Errorf("AtMost: incorrect horizontal lines for %v:\n   expected:%v\n   got:     %v", v.bounds, v.hlines, hlines)
		}
	}
}
