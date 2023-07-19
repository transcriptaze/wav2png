package wav2png

import (
	"testing"
)

func TestPCMToPixelScaling(t *testing.T) {
	tests := []struct {
		pcm      int
		height   int
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
		expected int16
	}{
		{0, 16, 0},
		{-1, 16, -1},
		{1, 16, 1},
		{-2, 16, -2},
		{32767, 16, 32767},
		{-32768, 16, -32768},
	}

	for _, test := range tests {
		if v := rescale(test.pcm, test.bits); v != test.expected {
			t.Errorf("Incorrectly rescaled value for %d-bit value %d: expected:%d, got:%d", test.bits, test.pcm, test.expected, v)
		}
	}

	for pcm := -32768; pcm < 32768; pcm++ {
		if v := rescale(pcm, 16); int(v) != pcm {
			t.Errorf("Incorrectly rescaled value for %d-bit value %d: expected:%d, got:%d", 16, pcm, pcm, v)
		}
	}
}

func TestVScale(t *testing.T) {
	tests := []struct {
		value    int16
		height   int
		expected int16
	}{
		// even height
		{+32767, 256, 255},
		{+256, 256, 129},
		{+255, 256, 128},
		{+1, 256, 128},
		{0, 256, 128},
		{-1, 256, 127},
		{-256, 256, 127},
		{-257, 256, 126},
		{-32768, 256, 0},

		// odd height
		{+32767, 257, 256},
		{+256, 257, 129},
		{+255, 257, 129},
		{+128, 257, 129},
		{+127, 257, 128},
		{+1, 257, 128},
		{0, 257, 128},
		{-1, 257, 128},
		{-127, 257, 128},
		{-128, 257, 127},
		{-256, 257, 127},
		{-257, 257, 127},
		{-32768, 257, 0},
	}

	for _, test := range tests {
		if h := vscale(test.value, test.height); h != test.expected {
			t.Errorf("%d incorrectly scaled to %d pixels: expected:%v, got:%v", test.value, test.height, test.expected, h)
		}
	}
}

func TestVScaleWithNegativeHeight(t *testing.T) {
	tests := []struct {
		value    int16
		height   int
		expected int16
	}{
		// even height
		{+32767, -256, 0},
		{+256, -256, 126},
		{+255, -256, 127},
		{+1, -256, 127},
		{0, -256, 127},
		{-1, -256, 128},
		{-256, -256, 128},
		{-257, -256, 129},
		{-32768, -256, 255},

		// odd height
		{+32767, -257, 0},
		{+256, -257, 127},
		{+255, -257, 127},
		{+128, -257, 127},
		{+127, -257, 128},
		{+1, -257, 128},
		{0, -257, 128},
		{-1, -257, 128},
		{-127, -257, 128},
		{-128, -257, 129},
		{-256, -257, 129},
		{-257, -257, 129},
		{-32768, -257, 256},
	}

	for _, test := range tests {
		if h := vscale(test.value, test.height); h != test.expected {
			t.Errorf("%d incorrectly scaled to %d pixels: expected:%v, got:%v", test.value, test.height, test.expected, h)
		}
	}
}
