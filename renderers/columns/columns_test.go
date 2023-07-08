package columns

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"os"
	"reflect"
	"testing"

	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/palettes"
)

//go:embed reference.wav
var audio []byte

//go:embed reference.png
var reference []byte

func TestRender(t *testing.T) {
	renderer := Columns{
		BarWidth:  16,
		BarGap:    1,
		Palette:   palettes.Fire,
		AntiAlias: kernels.Vertical,
		VScale:    1.0,
	}

	audio := read()
	samples := mix(audio, []int{1}...)[0:16000]

	if img, err := renderer.Render(samples, 640, 480, 0); err != nil {
		t.Fatalf("error rendering test image (%v)", err)
	} else if !reflect.DeepEqual(encode(img), reference) {
		t.Errorf("incorrectly rendered test image")

		os.WriteFile("../../runtime/columns.png", encode(img), 0666)
	}
}

func read() encoding.Audio {
	r := bytes.NewBuffer(audio)
	w, _ := wav.Decode(r)

	return encoding.Audio{
		SampleRate: float64(w.Format.SampleRate),
		Format:     fmt.Sprintf("%v", w.Format),
		Channels:   int(w.Format.Channels),
		Duration:   w.Duration(),
		Length:     w.Frames(),
		Samples:    w.Samples,
	}
}

func mix(wav encoding.Audio, channels ...int) []float32 {
	L := wav.Length
	N := float64(len(channels))
	samples := make([]float32, L)

	if len(wav.Samples) < 2 {
		return wav.Samples[0]
	}

	for i := 0; i < L; i++ {
		sample := 0.0
		for _, ch := range channels {
			sample += float64(wav.Samples[ch-1][i])
		}

		samples[i] = float32(sample / N)
	}

	return samples
}

func encode(img *image.NRGBA) []byte {
	var b bytes.Buffer
	png.Encode(&b, img)

	return b.Bytes()
}
