package wav2png

import (
	"image"
	"reflect"
	"testing"
)

func TestSquareGridSpecBorder(t *testing.T) {
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

func TestSquareGridSpecBorderWithPadding(t *testing.T) {
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

func TestSquareGridSpecVLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{64, 128, 192, 256, 320, 384, 448, 512, 576}
	gridspec := SquareGrid{size: 64, padding: 0}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridSpecVLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{65, 129, 193, 257, 321, 385, 449, 513, 577}
	gridspec := SquareGrid{size: 64, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridSpecVLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 645, 388)
	expected := []int{65, 129, 194, 258, 322, 386, 450, 515, 579}
	gridspec := SquareGrid{size: 64, padding: 1}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestSquareGridSpecHLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{192, 128, 64, 193, 257, 321}
	gridspec := SquareGrid{size: 64, padding: 0}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestSquareGridSpecHLinesWithPadding(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 388)
	expected := []int{193, 129, 65, 194, 258, 322}
	gridspec := SquareGrid{size: 64, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}

func TestSquareGridSpecHLinesWithNonIntegralSize(t *testing.T) {
	bounds := image.Rect(0, 0, 643, 390)
	expected := []int{194, 130, 65, 195, 260, 324}
	gridspec := SquareGrid{size: 64, padding: 1}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}
