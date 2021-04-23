package wav2png

import (
	"image"
	"reflect"
	"testing"
)

func TestNoGridBorder(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	gridspec := NoGrid{padding: 0}

	border := gridspec.Border(bounds)
	if border != nil {
		t.Errorf("NoGrid.Border should returnunexpectedly returned %v", border)
	}
}

func TestNoGridVLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{}
	gridspec := NoGrid{padding: 0}

	vlines := gridspec.VLines(bounds)
	if !reflect.DeepEqual(vlines, expected) {
		t.Errorf("Incorrect vertical lines:\n   expected:%v\n   got:     %v", expected, vlines)
	}
}

func TestNoGridHLines(t *testing.T) {
	bounds := image.Rect(0, 0, 641, 386)
	expected := []int{}
	gridspec := NoGrid{padding: 0}

	hlines := gridspec.HLines(bounds)
	if !reflect.DeepEqual(hlines, expected) {
		t.Errorf("Incorrect horizontal lines:\n   expected:%v\n   got:     %v", expected, hlines)
	}
}
