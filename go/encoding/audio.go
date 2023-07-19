package encoding

import (
	"fmt"
	"io"
	"time"

	"github.com/transcriptaze/wav2png/go/encoding/wav"
)

type Audio struct {
	SampleRate float64
	Format     string
	Channels   int
	Duration   time.Duration
	Length     int
	Samples    [][]float32
}

func Decode(f io.Reader) (audio Audio, err error) {
	var a *wav.WAV

	if a, err = wav.Decode(f); err != nil {
		return
	}

	audio = Audio{
		SampleRate: float64(a.Format.SampleRate),
		Format:     fmt.Sprintf("%v", a.Format),
		Channels:   int(a.Format.Channels),
		Duration:   a.Duration(),
		Length:     a.Frames(),
		Samples:    a.Samples,
	}

	return
}
