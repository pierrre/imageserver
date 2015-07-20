package image

import (
	"image"
	"io"

	"github.com/pierrre/imageserver"
)

func init() {
	RegisterEncoder("test", &testEncoder{})
}

var _ Encoder = &testEncoder{}

type testEncoder struct{}

func (enc *testEncoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return nil
}

func (enc *testEncoder) Change(params imageserver.Params) bool {
	return false
}

var testImageBounds = image.Rect(0, 0, 256, 256)

func NewTestImage() image.Image {
	return image.NewRGBA(testImageBounds)
}
