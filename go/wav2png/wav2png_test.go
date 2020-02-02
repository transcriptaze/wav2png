package wav2png

import (
	"testing"
)

func TestPCMToPixelScaling(t *testing.T) {
	tests := []struct {
		pcm      int
		height   uint
		expected int16
	}{
		{-32768, 256, 0},
		{32767, 256, 255},
		{-1, 256, 127},
		{0, 256, 128},
	}

	for _, test := range tests {
		v := rescale(test.pcm, 16)
		h := vscale(v, test.height)

		if h != test.expected {
			t.Errorf("%d incorrectly scaled to %d pixels: expected:%v, got:%v", test.pcm, test.height, test.expected, h)
		}
	}
}

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
	tests := []struct {
		value    int32
		height   uint
		expected int16
	}{
		{-65535, 256, 0},
		{+65535, 256, 255},
		{1, 256, 128},
		{-1, 256, 127},

		{511, 256, 128},
		// ???{512, 256, 129},

		{-511, 256, 127},
		{-512, 256, 126},
	}

	for _, test := range tests {
		if h := vscale(test.value, test.height); h != test.expected {
			t.Errorf("%d incorrectly scaled to %d pixels: expected:%v, got:%v", test.value, test.height, test.expected, h)
		}
	}
}
