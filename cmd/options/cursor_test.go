package options

import (
	"math"
	"testing"
	"time"
)

func TestLinearFn(t *testing.T) {
	frames := 150
	window := 1 * time.Second
	duration := 5 * time.Second

	tests := []struct {
		frame int
		x     float64
		shift float64
	}{
		{0, 0.0, 0.0},
		{1, 0.00666667, 0.0},
		{150, 1.0, 0.0},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))
		offset := seconds((duration - window).Seconds() * float64(v.frame) / float64(frames))

		x, shift := linear(at, offset, window, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}

		if math.Abs(v.shift-shift) > 0.000001 {
			t.Errorf("Invalid 'shift' for t '%v' - expected:%.3f, got:%.3f", at, v.shift, shift)
		}
	}
}

func TestCentreFn(t *testing.T) {
	frames := 150
	window := 1 * time.Second
	duration := 5 * time.Second

	tests := []struct {
		frame int
		x     float64
		shift float64
	}{
		{0, 0.5, 0.5},
		{1, 0.5, 0.466666667},
		{14, 0.5, 0.033333333},
		{15, 0.5, 0.0},
		{16, 0.5, 0.0},
		{134, 0.5, 0.0},
		{135, 0.5, 0.0},
		{136, 0.5, -0.033333333},
		{149, 0.5, -0.466666667},
		{150, 0.5, -0.5},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))
		offset := seconds((duration - window).Seconds() * float64(v.frame) / float64(frames))

		x, shift := centre(at, offset, window, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}

		if math.Abs(v.shift-shift) > 0.000001 {
			t.Errorf("Invalid 'shift' for t '%v' - expected:%.3f, got:%.3f", at, v.shift, shift)
		}
	}
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
