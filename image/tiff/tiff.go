// Package tiff provides a TIFF Encoder.
package tiff

import (
	"image"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"golang.org/x/image/tiff"
)

// Encoder encodes an Image to TIFF.
type Encoder struct {
}

var opts = &tiff.Options{
	Compression: tiff.Deflate,
	Predictor:   true,
}

// Encode implements Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return tiff.Encode(w, nim, opts)
}

// Change implements Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("tiff", &Encoder{})
}
