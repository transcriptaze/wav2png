package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func decodex(r io.Reader) (*WAV, error) {
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

	data := []float32{}
	switch format.Format {
	case 1:
		if audio, err := data16S(r); err != nil {
			return nil, err
		} else {
			data = audio
		}

	case 3:
		fact, err := decodeFact(r)
		if err != nil {
			return nil, err
		} else if fact == nil {
			return nil, fmt.Errorf("Invalid WAV 'fact ' subchunk")
		}

		if audio, err := dataF32(r); err != nil {
			return nil, err
		} else {
			data = audio
		}

	case 65534:
		if fmt.Sprintf("%0x", format.Extension.SubFormatGUID) != PCM_FLOAT {
			return nil, fmt.Errorf("Unsupported WAV file extension format %0x", format.Extension.SubFormatGUID)
		} else if format.BitsPerSample != 32 {
			return nil, fmt.Errorf("Unsupported sample format (float%v)", format.BitsPerSample)
		}

		if audio, err := dataWFX(r); err != nil {
			return nil, err
		} else {
			data = audio
		}

	default:
		return nil, fmt.Errorf("Unsupported WAV file format")
	}

	channels := int(format.Channels)
	samples := make([][]float32, channels)
	N := len(data) / channels
	for i := 0; i < channels; i++ {
		samples[i] = make([]float32, N)
	}

	frame := 0
	ix := 0
	for ix < len(data) {
		for i := 0; i < channels; i++ {
			samples[i][frame] = data[ix]
			ix++
		}
		frame++
	}

	return &WAV{
		Format:  *format,
		Samples: samples,
		frames:  frame,
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
	var extensionSize uint16
	var validBitsPerSample uint16
	var channelMask uint32
	var guid = make([]byte, 16)

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "fmt " {
		return nil, fmt.Errorf("Invalid 'fmt ' chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	} else if length != 16 && length != 40 {
		return nil, fmt.Errorf("Invalid 'fmt ' length %v - expected 16 (16-bit PCM) or 40 (32-bit floating point PCM)", length)
	}

	if err := binary.Read(r, binary.LittleEndian, &format); err != nil {
		return nil, err
	} else if format != 1 && format != 3 && format != 65534 {
		return nil, fmt.Errorf("Invalid 'fmt ' format %v - expected 1 (16-bit PCM), 3 (32-bit float PCM)  or 65534 (extensible)", format)
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
	} else if bitsPerSample != 16 && bitsPerSample != 32 {
		return nil, fmt.Errorf("Invalid 'fmt ' bits per sample %v - expected 16 (16-bit PCM) or 32 (32-bit PCM)", bitsPerSample)
	}

	if format == 65534 {
		if err := binary.Read(r, binary.LittleEndian, &extensionSize); err != nil {
			return nil, err
		} else if extensionSize != 22 {
			return nil, fmt.Errorf("Invalid extension size %v - expected 22 (extensible)", extensionSize)
		}

		if err := binary.Read(r, binary.LittleEndian, &validBitsPerSample); err != nil {
			return nil, err
		} else if validBitsPerSample != 32 {
			return nil, fmt.Errorf("Invalid 'valid bits per sample' extension field %v - expected 32 (32-bit floating point PCM)", validBitsPerSample)
		}

		if err := binary.Read(r, binary.LittleEndian, &channelMask); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.LittleEndian, &guid); err != nil {
			return nil, err
		}
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
		Extension: Extension{
			Length:             extensionSize,
			ValidBitsPerSample: validBitsPerSample,
			ChannelMask:        channelMask,
			SubFormatGUID:      guid,
		},
	}, nil
}

func decodeFact(r io.Reader) (*Fact, error) {
	var chunkID = make([]byte, 4)
	var length uint32
	var sampleFrames uint32

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "fact" {
		return nil, fmt.Errorf("Invalid 'fact ' chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	} else if length != 4 {
		return nil, fmt.Errorf("Invalid 'fact' length %v - expected 4", length)
	}

	if err := binary.Read(r, binary.LittleEndian, &sampleFrames); err != nil {
		return nil, err
	}

	return &Fact{
		ChunkID:      string(chunkID),
		Length:       length,
		SampleFrames: sampleFrames,
	}, nil
}

func data16S(r io.Reader) ([]float32, error) {
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

	return transcode(data)
}

func dataF32(r io.Reader) ([]float32, error) {
	var chunkID = make([]byte, 4)
	var length uint32

	for {
		if _, err := r.Read(chunkID); err != nil {
			return nil, err
		}

		if string(chunkID) != "data" {
			fmt.Printf("   discarding chunk ID '%s'\n", string(chunkID))

			if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
				return nil, err
			}

			chunk := make([]byte, length)
			if _, err := io.ReadFull(r, chunk); err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]float32, length/4)
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return data, nil
}

func dataWFX(r io.Reader) ([]float32, error) {
	var chunkID = make([]byte, 4)
	var length uint32
	var data []float32

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	} else if string(chunkID) != "data" {
		return nil, fmt.Errorf("Invalid 'data' chunk ID (%s)", string(chunkID))
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data = make([]float32, length/4)
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, err
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

		samples[i] = float32((2*int32(sample))+1) / 65536.0
	}

	return samples, nil
}
