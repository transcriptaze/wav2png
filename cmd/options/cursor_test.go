package options

import (
	"math"
	"testing"
	"time"
)

func TestLinearFn(t *testing.T) {
	cursor := Cursor{fn: "linear"}
	frames := 150
	duration := 5 * time.Second

	tests := []struct {
		frame  int
		offset time.Duration
		x      float64
		shift  float64
	}{
		{0, 0 * time.Second, 0.0, 0.0},
		{1, 27 * time.Millisecond, 0.00666667, 0.0},
		{30, 800 * time.Millisecond, 0.2, 0.0},
		{75, 2 * time.Second, 0.5, 0.0},
		{149, 3973 * time.Millisecond, 0.993333333, 0.0},
		{150, 4 * time.Second, 1.0, 0.0},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))

		x := cursor.Fn()(at, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}
	}
}

func TestCentreFn(t *testing.T) {
	cursor := Cursor{fn: "centre"}
	frames := 150
	duration := 5 * time.Second

	tests := []struct {
		frame  int
		offset time.Duration
		x      float64
		shift  float64
	}{
		{0, 0 * time.Second, 0.5, 0.5},
		{1, 0 * time.Second, 0.5, 0.466666667},
		{14, 0 * time.Second, 0.5, 0.033333333},
		{15, 0 * time.Second, 0.5, 0.0},
		{16, 33 * time.Millisecond, 0.5, 0.0},
		{75, 2 * time.Second, 0.5, 0.0},
		{134, 3967 * time.Millisecond, 0.5, 0.0},
		{135, 4 * time.Second, 0.5, 0.0},
		{136, 4 * time.Second, 0.5, -0.033333333},
		{149, 4 * time.Second, 0.5, -0.466666667},
		{150, 4 * time.Second, 0.5, -0.5},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))

		x := cursor.Fn()(at, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}
	}
}

func TestLeftFn(t *testing.T) {
	cursor := Cursor{fn: "left"}
	frames := 150
	duration := 5 * time.Second

	tests := []struct {
		frame  int
		offset time.Duration
		x      float64
		shift  float64
	}{
		{0, 0 * time.Second, 0.0, 0.0},
		{75, 2500 * time.Millisecond, 0.0, 0.0},
		{120, 4 * time.Second, 0.0, 0.0},
		{121, 4 * time.Second, 0.0, -0.033333333},
		{150, 4 * time.Second, 0.0, -1.0},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))

		x := cursor.Fn()(at, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}
	}
}

func TestRightFn(t *testing.T) {
	cursor := Cursor{fn: "right"}
	frames := 150
	duration := 5 * time.Second

	tests := []struct {
		frame  int
		offset time.Duration
		x      float64
		shift  float64
	}{
		{0, 0 * time.Second, 1.0, 1.0},
		{1, 0 * time.Second, 1.0, 0.966666667},
		{29, 0 * time.Millisecond, 1.0, 0.033333333},
		{30, 0 * time.Millisecond, 1.0, 0.0},
		{31, 33 * time.Millisecond, 1.0, 0.0},
		{75, 1500 * time.Millisecond, 1.0, 0.0},
		{150, 4 * time.Second, 1.0, 0.0},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))

		x := cursor.Fn()(at, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}
	}
}

func TestEaseFn(t *testing.T) {
	cursor := Cursor{fn: "ease"}
	frames := 150
	duration := 5 * time.Second

	tests := []struct {
		frame  int
		offset time.Duration
		x      float64
		shift  float64
	}{
		{0, 0 * time.Second, 0.0, 0.0},
		{1, 17 * time.Millisecond, 0.016666666, 0.0},
		{29, 483 * time.Millisecond, 0.483333333, 0.0},
		{30, 500 * time.Millisecond, 0.5, 0.0},
		{75, 2 * time.Second, 0.5, 0.0},
		{120, 3500 * time.Millisecond, 0.5, 0.0},
		{121, 3517 * time.Millisecond, 0.516666666, 0.0},
		{150, 4 * time.Second, 1.0, 0.0},
	}

	for _, v := range tests {
		at := seconds(duration.Seconds() * float64(v.frame) / float64(frames))

		x := cursor.Fn()(at, duration)

		if math.Abs(v.x-x) > 0.000001 {
			t.Errorf("Invalid cursor 'X' for t '%v' - expected:%.3f, got:%.3f", at, v.x, x)
		}
	}
}
