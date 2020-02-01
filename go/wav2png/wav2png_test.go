package wav2png

import (
	"testing"
)

func TestRescale(t *testing.T) {
	tests := []struct {
		pcm      int
		bits     uint
		expected int32
	}{
		{0, 16, 1},
		{-1, 16, -1},
		{1, 16, 3},
		{-2, 16, -3},
		{32767, 16, 65535},
		{-32768, 16, -65535},
	}

	for _, test := range tests {
		if v := rescale(test.pcm, test.bits); v != test.expected {
			t.Errorf("Incorrectly rescaled value for %d-bit value %d: expected:%d, got:%d", test.bits, test.pcm, test.expected, v)
		}
	}
}

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
