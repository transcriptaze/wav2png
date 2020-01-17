package wav2png

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestSinglePixelAntiAlias(t *testing.T) {
	vector := []struct {
		kernel [][]uint32
		colour color.NRGBA
	}{
		{none, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{vertical, color.NRGBA{R: 0xbf, G: 0xbf, B: 0xbf, A: 0xbf}},
		{horizontal, color.NRGBA{R: 0xaf, G: 0xaf, B: 0xaf, A: 0xaf}},
		{soft, color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xaa}},
	}

	for _, v := range vector {
		img := makeSinglePixelImage()
		expected := image.NewNRGBA(image.Rect(0, 0, 3, 3))

		expected.Set(1, 1, v.colour)

		result := antialias(img, v.kernel)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Result does not match expected:\n  expected:%v\n  got:     %v", expected, result)
		}
	}
}

func TestThreePixelNoAntiAlias(t *testing.T) {
	img := makeThreePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	expected.Set(1, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	expected.Set(1, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	result := antialias(img, none)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestThreePixelVerticalAntiAlias(t *testing.T) {
	img := makeThreePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(1, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(2, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(3, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})

	expected.Set(1, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	expected.Set(1, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(2, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(3, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})

	result := antialias(img, vertical)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestThreePixelHorizontalAntiAlias(t *testing.T) {
	img := makeThreePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(1, 1, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 1, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})

	expected.Set(1, 2, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})

	expected.Set(1, 3, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 3, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})

	result := antialias(img, horizontal)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestThreePixelSoftAntiAlias(t *testing.T) {
	img := makeThreePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(1, 1, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(2, 1, color.NRGBA{R: 0xe9, G: 0xe9, B: 0xe9, A: 0xe9})
	expected.Set(3, 1, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})

	expected.Set(1, 2, color.NRGBA{R: 0xe4, G: 0xe4, B: 0xe4, A: 0xe4})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xe4, G: 0xe4, B: 0xe4, A: 0xe4})

	expected.Set(1, 3, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(2, 3, color.NRGBA{R: 0xe9, G: 0xe9, B: 0xe9, A: 0xe9})
	expected.Set(3, 3, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})

	result := antialias(img, soft)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result does not match expected:\n  expected:%v\n  got:     %v", expected, result)
	}
}

func makeSinglePixelImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	img.Set(0, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	img.Set(0, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(0, 2, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	img.Set(1, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	img.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(1, 2, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})

	img.Set(2, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	img.Set(2, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(2, 2, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	return img
}

func makeThreePixelImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	img.Set(0, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	img.Set(0, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(0, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(0, 3, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(0, 4, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	img.Set(1, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	img.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(1, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(1, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(1, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})

	img.Set(2, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	img.Set(2, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(2, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(2, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})

	img.Set(3, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	img.Set(3, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(3, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(3, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(3, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})

	img.Set(4, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	img.Set(4, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(4, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(4, 3, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	img.Set(4, 4, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	return img
}

func diff(result, expected *image.NRGBA, t *testing.T) {
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 3; x++ {
			xy := image.Point{x, y}
			p := expected.At(xy.X, xy.Y)
			q := result.At(xy.X, xy.Y)

			if !reflect.DeepEqual(p, q) {
				t.Errorf("Pixel at %v does not match expected:\n  expected:%v\n  got:     %v", xy, p, q)
			}
		}
	}
}
