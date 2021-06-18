package wav

import (
	"math"
	"testing"
)

//func TestTranscode(t *testing.T) {
//	expected := []float32{
//		0.000000,
//
//		0.000031,
//		0.000244,
//		0.000488,
//		0.003906,
//		0.007813,
//		0.062502,
//		0.125004,
//		1.000000,
//
//		-0.000031,
//		-0.000244,
//		-0.000488,
//		-0.003906,
//		-0.007813,
//		-0.062502,
//		-0.125004,
//		-1.000031,
//	}
//
//	bytes := []byte{
//		0x00, 0x00,
//
//		0x01, 0x00,
//		0x08, 0x00,
//		0x10, 0x00,
//		0x80, 0x00,
//		0x00, 0x01,
//		0x00, 0x08,
//		0x00, 0x10,
//		0xff, 0x7f,
//
//		0xff, 0xff,
//		0xf8, 0xff,
//		0xf0, 0xff,
//		0x80, 0xff,
//		0x00, 0xff,
//		0x00, 0xf8,
//		0x00, 0xf0,
//		0x00, 0x80,
//	}
//
//	audio, err := transcode(bytes)
//	if err != nil {
//		t.Fatalf("Error transcoding valid data (%v)", err)
//	}
//
//	for i, v := range expected {
//		if math.Abs(float64(audio[i])-float64(v)) > 0.000001 {
//			t.Errorf("Incorrectly transcoded\n   expected:%.6f\n   got:     %.6f", expected, audio)
//			break
//		}
//	}
//}

func TestTranscode(t *testing.T) {
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

	audio, err := transcode(bytes)
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
