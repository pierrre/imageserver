// Package png provides a PNG Encoder.
package png

import (
	"image"
	"image/png"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// Encoder encodes an Image to PNG.
type Encoder struct {
	CompressionLevel png.CompressionLevel
}

// Encode implements Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	e := &png.Encoder{CompressionLevel: enc.CompressionLevel}
	return e.Encode(w, nim)
}

// Change implements Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("png", &Encoder{})
}
