package wav

import (
	"testing"
)

func TestInt24ToInt32(t *testing.T) {
	tests := []struct {
		i24      int24
		expected int32
	}{
		{int24{0x00, 0x00, 0x00}, 0},
		{int24{0x01, 0x00, 0x00}, 1},
		{int24{0xff, 0xff, 0x7f}, 8388607},
		{int24{0xff, 0xff, 0xff}, -1},
		{int24{0xfe, 0xff, 0xff}, -2},
		{int24{0x00, 0x00, 0x80}, -8388608},
	}

	for _, test := range tests {
		if v := test.i24.ToInt32(); v != test.expected {
			t.Errorf("invalid int24 value %v - expected:%v, got:%v", test.i24, test.expected, v)
		}
	}
}

func TestInt24ToFloat(t *testing.T) {
	tests := []struct {
		i24      int24
		expected float32
	}{
		{int24{0x00, 0x00, 0x00}, 1 / 16777216.0},
		{int24{0x01, 0x00, 0x00}, 3 / 16777216.0},
		{int24{0xff, 0xff, 0x7f}, 16777215 / 16777216.0},
		{int24{0xff, 0xff, 0xff}, -1 / 16777216.0},
		{int24{0xfe, 0xff, 0xff}, -3 / 16777216.0},
		{int24{0x00, 0x00, 0x80}, -16777215 / 16777216.0},
	}

	for _, test := range tests {
		if v := test.i24.ToFloat(); v != test.expected {
			t.Errorf("invalid int24 float value %v - expected:%v, got:%v", test.i24, test.expected, v)
		}
	}
}
