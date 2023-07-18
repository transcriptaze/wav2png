package compositor

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/transcriptaze/wav2png/encoding"
	"github.com/transcriptaze/wav2png/encoding/wav"
	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/palettes"
	"github.com/transcriptaze/wav2png/renderers/columns"
	"github.com/transcriptaze/wav2png/renderers/lines"
)

//go:embed test.wav
var audio []byte

//go:embed lines.png
var image_lines []byte

//go:embed columns.png
var image_columns []byte

var black = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
var green = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}

func TestLines(t *testing.T) {
	compositor := Compositor{
		width:      640,
		height:     480,
		padding:    0,
		scale:      1.0,
		background: fills.NewSolidFill(black),
		grid:       grids.NewSquareGrid(green, 64, grids.Approximate, false),

		renderer: lines.Lines{
			Palette:   palettes.Fire.Palette(),
			AntiAlias: kernels.Vertical,
		},
	}

	audio := read()
	from := 0 * time.Second
	to := 2 * time.Second
	fs := audio.SampleRate
	samples := mix(audio, []int{1}...)

	start := int(math.Floor(from.Seconds() * fs))
	end := int(math.Floor(to.Seconds() * fs))

	if img, err := compositor.Render(samples[start:end]); err != nil {
		t.Fatalf("error rendering test image (%v)", err)
	} else if !reflect.DeepEqual(encode(img), image_lines) {
		t.Errorf("incorrectly rendered lines image")

		os.WriteFile("../../runtime/columns.png", encode(img), 0666)
	}
}

func TestColumns(t *testing.T) {
	compositor := Compositor{
		width:      640,
		height:     480,
		padding:    0,
		background: fills.NewSolidFill(black),
		grid:       grids.NewSquareGrid(green, 64, grids.Approximate, false),

		renderer: columns.Columns{
			BarWidth:  16,
			BarGap:    1,
			Palette:   palettes.Fire.Palette(),
			AntiAlias: kernels.Vertical,
			VScale:    1.0,
		},
	}

	audio := read()
	from := 0 * time.Second
	to := 2 * time.Second
	fs := audio.SampleRate
	samples := mix(audio, []int{1}...)

	start := int(math.Floor(from.Seconds() * fs))
	end := int(math.Floor(to.Seconds() * fs))

	if img, err := compositor.Render(samples[start:end]); err != nil {
		t.Fatalf("error rendering test image (%v)", err)
	} else if !reflect.DeepEqual(encode(img), image_columns) {
		t.Errorf("incorrectly rendered columns image")

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
