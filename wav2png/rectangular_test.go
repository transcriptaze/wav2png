package wav2png

import (
	"image"
	"reflect"
	"testing"
)

func TestRectangularGridBorder(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := image.Rect(0, 0, 640, 385)
	gridspec := RectangularGrid{width: 64, height: 48, padding: 0}

	border := gridspec.Border(bounds)
	if border == nil {
		t.Fatalf("GridSpec.Border unexpectedly returned %v", border)
	}

	if !reflect.DeepEqual(*border, expected) {
		t.Errorf("Incorrect border:\n   expected:%v\n   got:     %v", expected, border)
	}
}

func TestRectangularGridBorderWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := image.Rect(1, 1, 641, 386)
	gridspec := RectangularGrid{width: 64, height: 48, padding: 1}

	border := gridspec.Border(bounds)
	if border == nil {
		t.Fatalf("GridSpec.Border unexpectedly returned %v", border)
	}

	if !reflect.DeepEqual(*border, expected) {
		t.Errorf("Incorrect border:\n   expected:%v\n   got:     %v", expected, border)
	}
}

func TestRectangularGridVLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{64, 128, 192, 256, 320, 384, 448, 512, 576}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 0}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestRectangularGridVLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{65, 129, 193, 257, 321, 385, 449, 513, 577}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestRectangularGridVLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 645, 388)
	expected := []int{65, 129, 194, 258, 322, 386, 450, 515, 579}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestRectangularGridHLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{192, 144, 96, 48, 193, 241, 289, 337}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 0}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestRectangularGridHLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{193, 145, 97, 49, 194, 242, 290, 338}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestRectangularGridHLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 390)
	expected := []int{194, 146, 97, 49, 195, 243, 292, 340}
	gridspec := RectangularGrid{width: 64, height: 48, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestRectangularGridApproximateFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 193, 241, 289, 337}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 383, 447, 511, 575}, hlines: []int{192, 144, 96, 48, 192, 240, 288, 336}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{63, 127, 190, 254, 317, 380, 444, 507, 571}, hlines: []int{189, 142, 94, 47, 190, 237, 285, 332}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{63, 126, 189, 252, 315, 378, 441, 504, 567}, hlines: []int{189, 142, 94, 47, 190, 237, 285, 332}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{65, 130, 195, 260, 325, 390, 455, 520, 585}, hlines: []int{195, 146, 97, 48, 196, 245, 294, 343}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 174, 124, 74, 24, 225, 275, 325, 375, 425}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 114, 68, 22, 161, 207, 253, 299}},
	}

	gridspec := RectangularGrid{width: 64, height: 48, padding: 0, fit: Approximate}

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

func TestRectangularGridExactFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 193, 241, 289, 337}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 192, 240, 288, 336}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 141, 93, 45, 190, 238, 286, 334}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 141, 93, 45, 190, 238, 286, 334}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{195, 147, 99, 51, 3, 196, 244, 292, 340, 388}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 112, 64, 16, 161, 209, 257, 305}},
	}

	gridspec := RectangularGrid{width: 64, height: 48, padding: 0, fit: Exact}

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

func TestRectangularGridAtLeastFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 193, 241, 289, 337}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 192, 240, 288, 336}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 141, 93, 45, 190, 238, 286, 334}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{189, 141, 93, 45, 190, 238, 286, 334}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{65, 130, 195, 260, 325, 390, 455, 520, 585}, hlines: []int{195, 146, 97, 48, 196, 245, 294, 343}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 174, 124, 74, 24, 225, 275, 325, 375, 425}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 112, 64, 16, 161, 209, 257, 305}},
	}

	gridspec := RectangularGrid{width: 64, height: 48, padding: 0, fit: AtLeast}

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

func TestRectangularGridAtMostFit(t *testing.T) {
	tests := []struct {
		bounds image.Rectangle
		vlines []int
		hlines []int
	}{
		{bounds: image.Rect(0, 0, 641, 386), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576}, hlines: []int{192, 144, 96, 48, 193, 241, 289, 337}},
		{bounds: image.Rect(0, 0, 640, 385), vlines: []int{64, 128, 192, 256, 320, 383, 447, 511, 575}, hlines: []int{192, 144, 96, 48, 192, 240, 288, 336}},
		{bounds: image.Rect(0, 0, 635, 380), vlines: []int{63, 127, 190, 254, 317, 380, 444, 507, 571}, hlines: []int{189, 142, 94, 47, 190, 237, 285, 332}},
		{bounds: image.Rect(0, 0, 631, 380), vlines: []int{63, 126, 189, 252, 315, 378, 441, 504, 567}, hlines: []int{189, 142, 94, 47, 190, 237, 285, 332}},
		{bounds: image.Rect(0, 0, 651, 392), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{195, 147, 99, 51, 3, 196, 244, 292, 340, 388}},
		{bounds: image.Rect(0, 0, 705, 450), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512, 576, 640}, hlines: []int{224, 176, 128, 80, 32, 225, 273, 321, 369, 417}},
		{bounds: image.Rect(0, 0, 577, 322), vlines: []int{64, 128, 192, 256, 320, 384, 448, 512}, hlines: []int{160, 114, 68, 22, 161, 207, 253, 299}},
	}

	gridspec := RectangularGrid{width: 64, height: 48, padding: 0, fit: AtMost}

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
