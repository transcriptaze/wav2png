package lines

import (
	_ "embed"
	"reflect"
	"testing"
)

//go:embed reference.png
var reference []byte

func TestRender(t *testing.T) {
	renderer := Lines{}

	if img, err := renderer.Render(); err != nil {
		t.Fatalf("error rendering test image (%v)", err)
	} else if !reflect.DeepEqual(img, reference) {
		t.Errorf("incorrectly rendered test image")
	}
}
