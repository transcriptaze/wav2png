package lines

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/styles/palettes"
	"github.com/transcriptaze/wav2png/wav2png"
)

//go:embed reference.wav
var audio []byte

//go:embed reference.png
var reference []byte

var black = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
var green = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}

func TestRender(t *testing.T) {
	renderer := Lines{
		Width:     640,
		Height:    480,
		Padding:   0,
		Palette:   palettes.Fire,
		FillSpec:  wav2png.NewSolidFill(black),
		GridSpec:  wav2png.NewSquareGrid(green, 64, wav2png.Approximate, false),
		AntiAlias: wav2png.Vertical,
		VScale:    1.0,
	}

	audio := read()
	from := 0 * time.Second
	to := 2 * time.Second
	fs := audio.SampleRate
	samples := mix(audio, []int{1}...)

	start := int(math.Floor(from.Seconds() * fs))
	end := int(math.Floor(to.Seconds() * fs))

	if img, err := renderer.Render(samples[start:end]); err != nil {
		t.Fatalf("error rendering test image (%v)", err)
	} else if !reflect.DeepEqual(encode(img), reference) {
		t.Errorf("incorrectly rendered test image")
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
