// Package bmp provides a BMP imageserver/image.Encoder implementation.
package bmp

import (
	"image"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"golang.org/x/image/bmp"
)

// Encoder is a BMP imageserver/image.Encoder implementation.
type Encoder struct{}

// Encode implements imageserver/image.Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return bmp.Encode(w, nim)
}

// Change implements imageserver/image.Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return false
}

func init() {
	imageserver_image.RegisterEncoder("bmp", &Encoder{})
}
