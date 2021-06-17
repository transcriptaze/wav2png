package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func Decode(r io.Reader) (*WAV, error) {
	header, err := decodeHeader(r)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, fmt.Errorf("Invalid WAV 'RIFF' chunk")
	}

	format, err := decodeFmt(r)
	if err != nil {
		return nil, err
	} else if format == nil {
		return nil, fmt.Errorf("Invalid WAV 'fmt ' subchunk")
	}

	data, err := decodeData(r)
	if err != nil {
		return nil, err
	}

	samples, err := transcode(data)
	if err != nil {
		return nil, err
	}

	return &WAV{
		Header:  *header,
		Format:  *format,
		Data:    data,
		Samples: samples,
	}, nil
}

func decodeHeader(r io.Reader) (*Header, error) {
	var chunkID = make([]byte, 4)
	var length uint32
	var format = make([]byte, 4)

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "RIFF" {
		return nil, fmt.Errorf("Invalid RIFF chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if _, err := r.Read(format); err != nil {
		return nil, err
	} else if string(format) != "WAVE" {
		return nil, fmt.Errorf("Invalid RIFF chunk format (%s)", string(format))
	}

	return &Header{
		ChunkID: string(chunkID),
		Length:  length,
		Format:  string(format),
	}, nil
}

func decodeFmt(r io.Reader) (*Format, error) {
	var chunkID = make([]byte, 4)
	var length uint32
	var format uint16
	var channels uint16
	var sampleRate uint32
	var byteRate uint32
	var blockAlign uint16
	var bitsPerSample uint16

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "fmt " {
		return nil, fmt.Errorf("Invalid 'fmt ' chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	} else if length != 16 {
		return nil, fmt.Errorf("Invalid 'fmt ' length %v - expected 16 (16-bit PCM)", length)
	}

	if err := binary.Read(r, binary.LittleEndian, &format); err != nil {
		return nil, err
	} else if format != 1 {
		return nil, fmt.Errorf("Invalid 'fmt ' format %v - expected 1 (16-bit PCM)", format)
	}

	if err := binary.Read(r, binary.LittleEndian, &channels); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &sampleRate); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &byteRate); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &blockAlign); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &bitsPerSample); err != nil {
		return nil, err
	} else if bitsPerSample != 16 {
		return nil, fmt.Errorf("Invalid 'fmt ' bits per sample %v - expected 16 (16-bit PCM)", format)
	}

	return &Format{
		ChunkID:       string(chunkID),
		Length:        length,
		Format:        format,
		Channels:      channels,
		SampleRate:    sampleRate,
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: bitsPerSample,
	}, nil
}

func decodeData(r io.Reader) ([]byte, error) {
	var chunkID = make([]byte, 4)
	var length uint32
	var data []byte

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "data" {
		return nil, fmt.Errorf("Invalid 'data' chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data = make([]byte, length)
	N, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	} else if N != int(length) {
		return nil, fmt.Errorf("Invalid 'data' length %v (expected %v)", N, length)
	}

	return data, nil
}

func transcode(data []byte) ([]float32, error) {
	N := len(data) / 2
	samples := make([]float32, N)
	r := bytes.NewReader(data)

	for i := 0; i < N; i++ {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return nil, err
		}

		samples[i] = float32(sample) / 32767.0
	}

	return samples, nil
}
