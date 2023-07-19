package wav

import (
	"bytes"
	_ "embed"
	"math"
	"reflect"
	"testing"
)

//go:embed PCM16.wav
var PCM16 []byte

//go:embed PCM24.wav
var PCM24 []byte

func TestDecodePCM16(t *testing.T) {
	expected := WAV{
		Format: Format{
			ChunkID:       "fmt ",
			Length:        16,
			Format:        1,
			Channels:      1,
			SampleRate:    8000,
			ByteRate:      16000,
			BlockAlign:    2,
			BitsPerSample: 16,
		},
		Samples: [][]float32{
			[]float32{0.02558899, 0.23469543, 0.45158386, 0.60517883, 0.69392395, 0.6950531, 0.61891174, 0.46516418, 0.26054382, 0.021255493},
		},

		frames: 100,
	}

	w, err := Decode(bytes.NewReader(PCM16))
	if err != nil {
		t.Fatalf("Error decoding WAV file (%v)", err)
	} else if w == nil {
		t.Fatalf("Failed to decode WAV file (%v)", w)
	}

	if !reflect.DeepEqual(w.Format, expected.Format) {
		t.Errorf("Invalid WAV 'fmt ' chunk \n   expected:%#v\n   got:     %#v", expected.Format, w.Format)
	}

	if len(w.Samples) != len(expected.Samples) {
		t.Errorf("Invalid WAV 'data' chunk \n   expected:%#v\n   got:     %#v", len(expected.Samples), len(w.Samples))
	} else if !reflect.DeepEqual(w.Samples[0][0:10], expected.Samples[0][0:10]) {
		t.Errorf("Invalid WAV 'data' chunk \n   expected:%#v\n   got:     %#v", expected.Samples[0][0:10], w.Samples[0][0:10])
	}

	if w.frames != expected.frames {
		t.Errorf("Invalid WAV 'data' chunk frames \n   expected:%#v\n   got:     %#v", expected.frames, w.frames)
	}
}

func TestDecodePCM24(t *testing.T) {
	expected := WAV{
		Format: Format{
			ChunkID:       "fmt ",
			Length:        40,
			Format:        WAVE_FORMAT_EXTENSIBLE,
			Channels:      1,
			SampleRate:    8000,
			ByteRate:      24000,
			BlockAlign:    3,
			BitsPerSample: 24,
			Extension: &Extension{
				Length:             22,
				ValidBitsPerSample: 24,
				ChannelMask:        0x04,
				SubFormatGUID:      []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71},
			},
		},
		Samples: [][]float32{
			[]float32{0.02557379, 0.23468024, 0.45156866, 0.60516363, 0.69390875, 0.6950379, 0.61889654, 0.465149, 0.26052862, 0.021240294},
		},

		frames: 100,
	}

	w, err := Decode(bytes.NewReader(PCM24))
	if err != nil {
		t.Fatalf("Error decoding WAV file (%v)", err)
	} else if w == nil {
		t.Fatalf("Failed to decode WAV file (%v)", w)
	}

	if w.Format.ChunkID != expected.Format.ChunkID ||
		w.Format.Length != expected.Format.Length ||
		w.Format.Format != expected.Format.Format ||
		w.Format.Channels != expected.Format.Channels ||
		w.Format.SampleRate != expected.Format.SampleRate ||
		w.Format.ByteRate != expected.Format.ByteRate ||
		w.Format.BlockAlign != expected.Format.BlockAlign ||
		w.Format.BitsPerSample != expected.Format.BitsPerSample {
		t.Errorf("Invalid WAV 'fmt ' chunk \n   expected:%#v\n   got:     %#v", expected.Format, w.Format)
	} else if w.Format.Extension == nil || !reflect.DeepEqual(*w.Format.Extension, *expected.Format.Extension) {
		t.Errorf("Invalid WAV 'fmt ' chunk \n   expected:%#v\n   got:     %#v", expected.Format.Extension, w.Format.Extension)
	}

	if len(w.Samples) != len(expected.Samples) {
		t.Errorf("Invalid WAV 'data' chunk \n   expected:%#v\n   got:     %#v", len(expected.Samples), len(w.Samples))
	} else if !reflect.DeepEqual(w.Samples[0][0:10], expected.Samples[0][0:10]) {
		t.Errorf("Invalid WAV 'data' chunk \n   expected:%#v\n   got:     %#v", expected.Samples[0][0:10], w.Samples[0][0:10])
	}

	if w.frames != expected.frames {
		t.Errorf("Invalid WAV 'data' chunk frames \n   expected:%#v\n   got:     %#v", expected.frames, w.frames)
	}
}

func TestPCM16(t *testing.T) {
	expected := []float32{
		0.000015,
		0.000046,
		0.000259,
		0.000504,
		0.003922,
		0.007828,
		0.062515,
		0.125015,
		0.999985,

		-0.000015,
		-0.000046,
		-0.000259,
		-0.000504,
		-0.003922,
		-0.062515,
		-0.125015,
		-0.999985,
	}

	bytes := []byte{
		0x00, 0x00,
		0x01, 0x00,
		0x08, 0x00,
		0x10, 0x00,
		0x80, 0x00,
		0x00, 0x01,
		0x00, 0x08,
		0x00, 0x10,
		0xff, 0x7f,

		0xff, 0xff,
		0xfe, 0xff,
		0xf7, 0xff,
		0xef, 0xff,
		0x7f, 0xff,
		0xff, 0xf7,
		0xff, 0xef,
		0x00, 0x80,
	}

	audio, err := parsePCM16(bytes)
	if err != nil {
		t.Fatalf("Error transcoding valid data (%v)", err)
	}

	for i, v := range expected {
		if math.Abs(float64(audio[i])-float64(v)) > 0.000001 {
			t.Errorf("Incorrectly transcoded\n   expected:%.6f\n   got:     %.6f", expected, audio)
			break
		}
	}
}
