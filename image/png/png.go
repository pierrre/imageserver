// Package png provides a PNG imageserver/image.Encoder implementation.
package png

import (
	"image"
	"image/png"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// Encoder is a PNG imageserver/image.Encoder implementation.
type Encoder struct {
	CompressionLevel png.CompressionLevel
}

// Encode implements imageserver/image.Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	e := &png.Encoder{CompressionLevel: enc.CompressionLevel}
	return e.Encode(w, nim)
}

// Change implements imageserver/image.Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("png", &Encoder{})
}
