package lines

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
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
	audio := read()
	from := 0 * time.Second
	to := 2 * time.Second

	renderer := Lines{
		Width:     640,
		Height:    480,
		Padding:   0,
		Palette:   palettes.Fire,
		FillSpec:  wav2png.NewSolidFill(black),
		GridSpec:  wav2png.NewSquareGrid(green, 64, wav2png.Approximate, false),
		AntiAlias: wav2png.Vertical,
		VScale:    1.0,
		Channels:  []int{1},
	}

	if img, err := renderer.Render(audio, from, to); err != nil {
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

func encode(img *image.NRGBA) []byte {
	var b bytes.Buffer
	png.Encode(&b, img)

	return b.Bytes()
}
