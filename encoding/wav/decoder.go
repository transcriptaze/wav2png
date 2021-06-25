package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type chunk struct {
	ID     string
	length uint32
	data   []byte
}

func Decode(r io.Reader) (*WAV, error) {
	var b *bytes.Buffer
	var w = WAV{}

	// ... get WAV header
	if chunk, err := getChunk(r); err != nil {
		return nil, err
	} else if chunk == nil {
		return nil, fmt.Errorf("Invalid RIFF header chunk (%v)", chunk)
	} else if chunk.ID != "RIFF" {
		return nil, fmt.Errorf("Invalid RIFF header chunk ID (%s)", chunk.ID)
	} else {
		b = bytes.NewBuffer(chunk.data)
	}

	var format = make([]byte, 4)
	if _, err := b.Read(format); err != nil {
		return nil, err
	} else if string(format) != "WAVE" {
		return nil, fmt.Errorf("Invalid WAV header format (%s)", string(format))
	}

	// ... read remaining chunks
	chunks := map[string]*chunk{}

	for {
		chunk, err := getChunk(b)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		chunks[chunk.ID] = chunk
	}

	// ... extract 'fmt '
	if chunk, ok := chunks["fmt "]; !ok {
		return nil, fmt.Errorf("Invalid WAV file - missing 'fmt ' subchunk")
	} else if f, err := parseFMT(*chunk); err != nil {
		return nil, fmt.Errorf("Invalid WAV 'fmt ' subchunk (%v)", err)
	} else if f == nil {
		return nil, fmt.Errorf("Invalid WAV 'fmt ' subchunk (%v)", f)
	} else {
		w.Format = *f
	}

	// ... extract 'fact'

	if chunk, ok := chunks["fact"]; ok {
		if fact, err := parseFact(*chunk); err != nil {
			return nil, fmt.Errorf("Invalid WAV file 'fact' subchunk (%v)", err)
		} else {
			w.Fact = fact
		}
	}

	// ... extract 'data'
	data := []float32{}
	if chunk, ok := chunks["data"]; !ok {
		return nil, fmt.Errorf("Invalid WAV file - missing 'data' subchunk")
	} else {
		switch w.Format.Format {
		case 1:
			if audio, err := parsePCM16(chunk.data); err != nil {
				return nil, err
			} else {
				data = audio
			}

		case 3:
			//		fact, err := decodeFact(r)
			//		if err != nil {
			//			return nil, err
			//		} else if fact == nil {
			//			return nil, fmt.Errorf("Invalid WAV 'fact ' subchunk")
			//		}
			if audio, err := parsePCM32f(chunk.data); err != nil {
				return nil, err
			} else {
				data = audio
			}

		case 65534:
			if fmt.Sprintf("%0x", w.Format.Extension.SubFormatGUID) != PCM_FLOAT {
				return nil, fmt.Errorf("Unsupported WAV file extension format %0x", w.Format.Extension.SubFormatGUID)
			} else if w.Format.BitsPerSample != 32 {
				return nil, fmt.Errorf("Unsupported sample format (float%v)", w.Format.BitsPerSample)
			}

			if audio, err := parseWFX(chunk.data); err != nil {
				return nil, err
			} else {
				data = audio
			}

		default:
			return nil, fmt.Errorf("Unsupported WAV file format")
		}
	}

	// ... split raw data into channels
	channels := int(w.Format.Channels)
	samples := make([][]float32, w.Format.Channels)
	N := len(data) / channels
	for i := 0; i < channels; i++ {
		samples[i] = make([]float32, N)
	}

	frames := 0
	ix := 0
	for ix < len(data) {
		for i := 0; i < channels; i++ {
			samples[i][frames] = data[ix]
			ix++
		}
		frames++
	}

	w.Samples = samples
	w.frames = frames

	return &w, nil
}

func getChunk(r io.Reader) (*chunk, error) {
	var chunkID = make([]byte, 4)
	var length uint32

	if _, err := r.Read(chunkID); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("Error reading chunk '%s' from WAV file (%s)", string(chunkID), err)
	}

	return &chunk{
		ID:     string(chunkID),
		length: length,
		data:   data,
	}, nil
}

func parseFMT(ch chunk) (*Format, error) {
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

	r := bytes.NewBuffer(ch.data)

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
		ChunkID:       ch.ID,
		Length:        ch.length,
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

func parseFact(ch chunk) (*Fact, error) {
	var sampleFrames uint32

	r := bytes.NewBuffer(ch.data)
	if err := binary.Read(r, binary.LittleEndian, &sampleFrames); err != nil {
		return nil, err
	}

	return &Fact{
		ChunkID:      ch.ID,
		Length:       ch.length,
		SampleFrames: sampleFrames,
	}, nil
}

func parsePCM16(b []byte) ([]float32, error) {
	return transcode(b)
}

func parsePCM32f(b []byte) ([]float32, error) {
	data := make([]float32, len(b)/4)
	r := bytes.NewBuffer(b)
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return data, nil
}

func parseWFX(b []byte) ([]float32, error) {
	data := make([]float32, len(b)/4)
	r := bytes.NewBuffer(b)
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return data, nil
}
