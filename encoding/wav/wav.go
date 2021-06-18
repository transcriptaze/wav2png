package wav

import (
	"fmt"
	"time"
)

type WAV struct {
	Header  Header
	Format  Format
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
	Extension     Extension
}

type Extension struct {
	Length             uint16
	ValidBitsPerSample uint16
	ChannelMask        uint32
	SubFormatGUID      []byte
}

type Data struct {
	ChunkID string
	Length  uint32
	Audio   []byte
}

func (w *WAV) Duration() time.Duration {
	d := float64(len(w.Samples)) / float64(w.Format.SampleRate)

	return time.Duration(d * float64(time.Second))
}

func (f Format) String() string {
	switch f.Format {
	case 1:
		if f.BitsPerSample == 16 {
			return "16-bit signed PCM"
		}

	case 65534:
		if fmt.Sprintf("%0x", f.Extension.SubFormatGUID) == PCM_FLOAT {
			if f.BitsPerSample == 32 {
				return "32-bit floating point PCM"
			}
		}
	}

	return "unknown"
}
