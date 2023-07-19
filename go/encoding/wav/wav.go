package wav

import (
	"fmt"
	"time"
)

const WAVE_FORMAT_PCM uint16 = 0x0001
const WAVE_FORMAT_IEEE_FLOAT uint16 = 0x0003
const WAVE_FORMAT_ALAW uint16 = 0x0006
const WAVE_FORMAT_MULAW uint16 = 0x0007
const WAVE_FORMAT_EXTENSIBLE uint16 = 0xFFFE

const GUID_PCM = "0100000000001000800000aa00389b71"
const GUID_IEEE_FLOAT = "0300000000001000800000aa00389b71"

type WAV struct {
	Format  Format
	Fact    *Fact
	Samples [][]float32
	frames  int
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
	Extension     *Extension
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
	return time.Duration(float64(w.frames) * float64(time.Second) / float64(w.Format.SampleRate))
}

func (f Format) String() string {
	switch {
	case f.Format == WAVE_FORMAT_PCM && f.BitsPerSample == 16:
		return "16-bit signed PCM"

	case f.Format == WAVE_FORMAT_IEEE_FLOAT && f.BitsPerSample == 32:
		return "32-bit floating point PCM"

	case f.Format == WAVE_FORMAT_EXTENSIBLE && fmt.Sprintf("%0x", f.Extension.SubFormatGUID) == GUID_PCM && f.BitsPerSample == 24:
		return "24-bit floating point PCM"

	case f.Format == WAVE_FORMAT_EXTENSIBLE && fmt.Sprintf("%0x", f.Extension.SubFormatGUID) == GUID_IEEE_FLOAT && f.BitsPerSample == 32:
		return "32-bit floating point PCM"
	}

	return "unknown"
}
