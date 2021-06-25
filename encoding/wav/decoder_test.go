package wav

import (
	"bytes"
	_ "embed"
	"reflect"
	"testing"
)

//go:embed PCM16.wav
var PCM16 []byte

func TestDecode(t *testing.T) {
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
			Extension: Extension{
				SubFormatGUID: []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
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
