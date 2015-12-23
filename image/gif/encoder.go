package gif

import (
	"image"
	"image/gif"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// Encoder encodes an Image to GIF.
type Encoder struct{}

// Encode implements Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return gif.Encode(w, nim, nil)
}

// Change implements Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("gif", &Encoder{})
}
