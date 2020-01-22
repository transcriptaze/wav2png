package wav2png

import (
	"testing"
)

func TestVScale(t *testing.T) {
	height := uint(256)
	vector := []struct {
		value    int
		expected int
	}{
		{-32768, 0},
		{-1, 127},
		{0, 128},
		{+1, 128},
		{32767, 255},
	}

	for _, v := range vector {
		if h := vscale(v.value, height); h != v.expected {
			t.Errorf("Incorrect scale for %v: expected:%v, got:%v", v.value, v.expected, h)
		}
	}
}
