package wav

import (
	"fmt"
	"time"
)

const PCM_FLOAT = "0300000000001000800000aa00389b71"

type WAV struct {
	Header  Header
	Format  Format
	Fact    *Fact
	Samples [][]float32
	frames  int
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

type Fact struct {
	ChunkID      string
	Length       uint32
	SampleFrames uint32
}

type Data struct {
	ChunkID string
	Length  uint32
	Audio   []byte
}

func (w *WAV) Frames() int {
	return w.frames
}

func (w *WAV) Duration() time.Duration {
	if w.frames > 0 {
		return time.Duration(float64(w.frames) * float64(time.Second) / float64(w.Format.SampleRate))
	}

	d := float64(len(w.Samples)) / float64(w.Format.SampleRate)

	return time.Duration(d * float64(time.Second))
}

func (f Format) String() string {
	switch f.Format {
	case 1:
		if f.BitsPerSample == 16 {
			return "16-bit signed PCM"
		}

	case 3:
		if f.BitsPerSample == 32 {
			return "32-bit floating point PCM"
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
