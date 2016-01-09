package gif

import (
	"image"
	"image/gif"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// Encoder is a GIF imageserver/image.Encoder implementation.
type Encoder struct{}

// Encode implements imageserver/image.Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return gif.Encode(w, nim, nil)
}

// Change implements imageserver/image.Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("gif", &Encoder{})
}
