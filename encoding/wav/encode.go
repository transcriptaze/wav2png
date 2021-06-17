package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func Encode(left, right []float32) ([]byte, error) {
	audio, err := transcodeF(left, right)
	if err != nil {
		return nil, err
	}

	data := Data{
		ChunkID: "data",
		Length:  uint32(len(audio)),
		Audio:   audio,
	}

	format := Format{
		ChunkID:       "fmt ",
		Length:        16,
		Format:        1,
		Channels:      2,
		SampleRate:    44100,
		ByteRate:      2 * 44100,
		BlockAlign:    2,
		BitsPerSample: 16,
	}

	header := Header{
		ChunkID: "RIFF",
		Length:  uint32(4 + (8 + int(format.Length)) + (8 + int(data.Length))),
		Format:  "WAVE",
	}

	var w bytes.Buffer

	if err := header.encode(&w); err != nil {
		return nil, err
	}

	if err := format.encode(&w); err != nil {
		return nil, err
	}

	if err := data.encode(&w); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (chunk *Header) encode(w io.Writer) error {
	if N, err := w.Write([]byte(chunk.ChunkID)); err != nil {
		return err
	} else if N != len(chunk.ChunkID) {
		return fmt.Errorf("Failed to write RIFF chunk ID")
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Length); err != nil {
		return err
	}

	if N, err := w.Write([]byte(chunk.Format)); err != nil {
		return err
	} else if N != len(chunk.Format) {
		return fmt.Errorf("Failed to write RIFF chunk format")
	}

	return nil
}

func (chunk *Format) encode(w io.Writer) error {
	if N, err := w.Write([]byte(chunk.ChunkID)); err != nil {
		return err
	} else if N != len(chunk.ChunkID) {
		return fmt.Errorf("Failed to write 'fmt' chunk ID")
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Length); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Format); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Channels); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.SampleRate); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.ByteRate); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.BlockAlign); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.BitsPerSample); err != nil {
		return err
	}

	return nil
}

func (chunk *Data) encode(w io.Writer) error {
	if N, err := w.Write([]byte(chunk.ChunkID)); err != nil {
		return err
	} else if N != len(chunk.ChunkID) {
		return fmt.Errorf("Failed to write 'data' chunk ID")
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Length); err != nil {
		return err
	}

	if N, err := w.Write(chunk.Audio); err != nil {
		return err
	} else if N != len(chunk.Audio) {
		return fmt.Errorf("Failed to write 'data' chunk audio")
	}

	return nil
}

func transcodeF(left, right []float32) ([]byte, error) {
	var b bytes.Buffer

	N := len(left)
	if N < len(right) {
		N = len(right)
	}

	for i := 0; i < N; i++ {
		var l float32
		var r float32

		if i < len(left) {
			l = left[i]
		}

		if i < len(right) {
			r = right[i]
		}

		if err := binary.Write(&b, binary.LittleEndian, clip(l)); err != nil {
			return nil, err
		}

		if err := binary.Write(&b, binary.LittleEndian, clip(r)); err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

func clip(v float32) int16 {
	vf := mixLogarithmicRangeCompression(v)
	sample := int64(vf * 32767.0)

	if sample > 32767 {
		return 32767
	}

	if sample < -32768 {
		return -32768
	}

	return int16(sample)
}

// Ref. https://github.com/go-mix/mix
func mixLogarithmicRangeCompression(v float32) float32 {
	if v < -1 {
		return float32(-math.Log(-float64(v)-0.85)/14 - 0.75)
	}

	if v > 1 {
		return float32(math.Log(float64(v)-0.85)/14 + 0.75)
	}

	return float32(v / 1.61803398875)
}
