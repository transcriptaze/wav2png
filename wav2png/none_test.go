package wav2png

import (
	"image"
	"reflect"
	"testing"
)

func TestNoGridBorder(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	padding := 0
	var gridspec GridSpec = NoGrid{}

	border := gridspec.Border(bounds, padding)
	if border != nil {
		t.Errorf("NoGrid.Border unexpectedly returned %v", border)
	}
}

func TestNoGridVLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	padding := 0
	expected := []int{}
	var gridspec GridSpec = NoGrid{}

	vlines := gridspec.VLines(bounds, padding)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestNoGridHLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	padding := 0
	expected := []int{}
	var gridspec GridSpec = NoGrid{}

	hlines := gridspec.HLines(bounds, padding)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}
