package encoding

import (
	"time"
)

type Audio struct {
	SampleRate float64
	Format     string
	Channels   int
	Duration   time.Duration
	Length     int
	Samples    [][]float32
}
