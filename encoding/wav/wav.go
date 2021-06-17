package wav

import (
	"time"
)

type WAV struct {
	Header  Header
	Format  Format
	Data    []byte
	Samples []float32
}

type Header struct {
	ChunkID string
	Length  uint32
	Format  string
}

type Format struct {
	ChunkID       string
	Length        uint32
	Format        uint16
	Channels      uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type Data struct {
	ChunkID string
	Length  uint32
	Audio   []byte
}

func (w *WAV) Duration() time.Duration {
	d := float64(len(w.Data)) / float64(w.Format.ByteRate)

	return time.Duration(d * float64(time.Second))
}
