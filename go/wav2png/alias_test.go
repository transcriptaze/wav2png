package wav2png

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestSinglePixelNoAntiAlias(t *testing.T) {
	img := makeSinglePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	expected.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	result := antialias(img, none)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestSinglePixelVerticalAntiAlias(t *testing.T) {
	img := makeSinglePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	expected.Set(1, 0, color.NRGBA{R: 0x3f, G: 0x3f, B: 0x3f, A: 0x3f})
	expected.Set(1, 1, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(1, 2, color.NRGBA{R: 0x3f, G: 0x3f, B: 0x3f, A: 0x3f})

	result := antialias(img, vertical)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestSinglePixelHorizontalAntiAlias(t *testing.T) {
	img := makeSinglePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	expected.Set(0, 1, color.NRGBA{R: 0x3f, G: 0x3f, B: 0x3f, A: 0x3f})
	expected.Set(1, 1, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(2, 1, color.NRGBA{R: 0x3f, G: 0x3f, B: 0x3f, A: 0x3f})

	result := antialias(img, horizontal)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestSinglePixelSoftAntiAlias(t *testing.T) {
	img := makeSinglePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	expected.Set(0, 0, color.NRGBA{R: 0x0a, G: 0x0a, B: 0x0a, A: 0x0a})
	expected.Set(1, 0, color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0x15})
	expected.Set(2, 0, color.NRGBA{R: 0x0a, G: 0x0a, B: 0x0a, A: 0x0a})

	expected.Set(0, 1, color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0x15})
	expected.Set(1, 1, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(2, 1, color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0x15})

	expected.Set(0, 2, color.NRGBA{R: 0x0a, G: 0x0a, B: 0x0a, A: 0x0a})
	expected.Set(1, 2, color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0x15})
	expected.Set(2, 2, color.NRGBA{R: 0x0a, G: 0x0a, B: 0x0a, A: 0x0a})

	result := antialias(img, soft)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestFivePixelNoAntiAlias(t *testing.T) {
	img := makeFivePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(0, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	expected.Set(1, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(2, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(3, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(4, 0, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	expected.Set(0, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	expected.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(4, 1, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})

	expected.Set(0, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	expected.Set(1, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(4, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})

	expected.Set(0, 3, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	expected.Set(1, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(4, 3, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})

	expected.Set(0, 4, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})
	expected.Set(1, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(2, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(3, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(4, 4, color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0x20})

	result := antialias(img, none)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestFivePixelVerticalAntiAlias(t *testing.T) {
	img := makeFivePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(0, 0, color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0x28})
	expected.Set(1, 0, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(2, 0, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(3, 0, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(4, 0, color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0x28})

	expected.Set(0, 1, color.NRGBA{R: 0x50, G: 0x50, B: 0x50, A: 0x50})
	expected.Set(1, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(2, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(3, 1, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(4, 1, color.NRGBA{R: 0x50, G: 0x50, B: 0x50, A: 0x50})

	expected.Set(0, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})
	expected.Set(1, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(4, 2, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0x60})

	expected.Set(0, 3, color.NRGBA{R: 0x50, G: 0x50, B: 0x50, A: 0x50})
	expected.Set(1, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(2, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(3, 3, color.NRGBA{R: 0xdf, G: 0xdf, B: 0xdf, A: 0xdf})
	expected.Set(4, 3, color.NRGBA{R: 0x50, G: 0x50, B: 0x50, A: 0x50})

	expected.Set(0, 4, color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0x28})
	expected.Set(1, 4, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(2, 4, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(3, 4, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(4, 4, color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0x28})

	result := antialias(img, vertical)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestFivePixelHorizontalAntiAlias(t *testing.T) {
	img := makeFivePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(0, 0, color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0x30})
	expected.Set(1, 0, color.NRGBA{R: 0x68, G: 0x68, B: 0x68, A: 0x68})
	expected.Set(2, 0, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(3, 0, color.NRGBA{R: 0x68, G: 0x68, B: 0x68, A: 0x68})
	expected.Set(4, 0, color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0x30})

	expected.Set(0, 1, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})
	expected.Set(1, 1, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 1, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(4, 1, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})

	expected.Set(0, 2, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})
	expected.Set(1, 2, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(4, 2, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})

	expected.Set(0, 3, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})
	expected.Set(1, 3, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(2, 3, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 3, color.NRGBA{R: 0xd7, G: 0xd7, B: 0xd7, A: 0xd7})
	expected.Set(4, 3, color.NRGBA{R: 0x6f, G: 0x6f, B: 0x6f, A: 0x6f})

	expected.Set(0, 4, color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0x30})
	expected.Set(1, 4, color.NRGBA{R: 0x68, G: 0x68, B: 0x68, A: 0x68})
	expected.Set(2, 4, color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0x80})
	expected.Set(3, 4, color.NRGBA{R: 0x68, G: 0x68, B: 0x68, A: 0x68})
	expected.Set(4, 4, color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0x30})

	result := antialias(img, horizontal)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func TestFivePixelSoftAntiAlias(t *testing.T) {
	img := makeFivePixelImage()
	expected := image.NewNRGBA(image.Rect(0, 0, 5, 5))

	expected.Set(0, 0, color.NRGBA{R: 0x2d, G: 0x2d, B: 0x2d, A: 0x2d})
	expected.Set(1, 0, color.NRGBA{R: 0x71, G: 0x71, B: 0x71, A: 0x71})
	expected.Set(2, 0, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(3, 0, color.NRGBA{R: 0x71, G: 0x71, B: 0x71, A: 0x71})
	expected.Set(4, 0, color.NRGBA{R: 0x2d, G: 0x2d, B: 0x2d, A: 0x2d})

	expected.Set(0, 1, color.NRGBA{R: 0x5f, G: 0x5f, B: 0x5f, A: 0x5f})
	expected.Set(1, 1, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(2, 1, color.NRGBA{R: 0xe9, G: 0xe9, B: 0xe9, A: 0xe9})
	expected.Set(3, 1, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(4, 1, color.NRGBA{R: 0x5f, G: 0x5f, B: 0x5f, A: 0x5f})

	expected.Set(0, 2, color.NRGBA{R: 0x6a, G: 0x6a, B: 0x6a, A: 0x6a})
	expected.Set(1, 2, color.NRGBA{R: 0xe4, G: 0xe4, B: 0xe4, A: 0xe4})
	expected.Set(2, 2, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	expected.Set(3, 2, color.NRGBA{R: 0xe4, G: 0xe4, B: 0xe4, A: 0xe4})
	expected.Set(4, 2, color.NRGBA{R: 0x6a, G: 0x6a, B: 0x6a, A: 0x6a})

	expected.Set(0, 3, color.NRGBA{R: 0x5f, G: 0x5f, B: 0x5f, A: 0x5f})
	expected.Set(1, 3, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(2, 3, color.NRGBA{R: 0xe9, G: 0xe9, B: 0xe9, A: 0xe9})
	expected.Set(3, 3, color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 0xd1})
	expected.Set(4, 3, color.NRGBA{R: 0x5f, G: 0x5f, B: 0x5f, A: 0x5f})

	expected.Set(0, 4, color.NRGBA{R: 0x2d, G: 0x2d, B: 0x2d, A: 0x2d})
	expected.Set(1, 4, color.NRGBA{R: 0x71, G: 0x71, B: 0x71, A: 0x71})
	expected.Set(2, 4, color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0x7f})
	expected.Set(3, 4, color.NRGBA{R: 0x71, G: 0x71, B: 0x71, A: 0x71})
	expected.Set(4, 4, color.NRGBA{R: 0x2d, G: 0x2d, B: 0x2d, A: 0x2d})

	result := antialias(img, soft)

	if !reflect.DeepEqual(result, expected) {
		diff(result, expected, t)
	}
}

func makeSinglePixelImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	img.Set(0, 0, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	img.Set(0, 1, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	img.Set(0, 2, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})

	img.Set(1, 0, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	img.Set(1, 1, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	img.Set(1, 2, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})

	img.Set(2, 0, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	img.Set(2, 1, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	img.Set(2, 2, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})

	return img
}

func makeFivePixelImage() *image.NRGBA {
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
	tl := result.Bounds().Min
	br := result.Bounds().Max

	for y := tl.Y; y < br.Y; y++ {
		for x := tl.X; x < br.X; x++ {
			xy := image.Point{x, y}
			p := expected.At(xy.X, xy.Y)
			q := result.At(xy.X, xy.Y)

			if !reflect.DeepEqual(p, q) {
				t.Errorf("Pixel at %v does not match expected:\n  expected:%v\n  got:     %v", xy, p, q)
			}
		}
	}
}
