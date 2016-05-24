// Package jpeg provides a JPEG imageserver/image.Encoder implementation.
package jpeg

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// Encoder is a JPEG imageserver/image.Encoder implementation.
//
// It supports the "quality" param (1 to 100).
type Encoder struct {
	DefaultQuality int
}

// Encode implements imageserver/image.Encoder.
func (enc *Encoder) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	opts, err := enc.getOptions(params)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, nim, opts)
}

func (enc *Encoder) getOptions(params imageserver.Params) (*jpeg.Options, error) {
	opts := &jpeg.Options{}
	var err error
	opts.Quality, err = enc.getQuality(params)
	if err != nil {
		return nil, err
	}
	return opts, nil
}

func (enc *Encoder) getQuality(params imageserver.Params) (int, error) {
	if !params.Has("quality") {
		if enc.DefaultQuality != 0 {
			return enc.DefaultQuality, nil
		}
		return jpeg.DefaultQuality, nil
	}
	quality, err := params.GetInt("quality")
	if err != nil {
		return 0, err
	}
	if quality < 1 {
		return 0, &imageserver.ParamError{Param: "quality", Message: "must be greater than or equal to 1"}
	}
	if quality > 100 {
		return 0, &imageserver.ParamError{Param: "quality", Message: "must be less than or equal to 100"}
	}
	return quality, nil
}

// Change implements imageserver/image.Encoder.
func (enc *Encoder) Change(params imageserver.Params) bool {
	return params.Has("quality")
}

func init() {
	imageserver_image.RegisterEncoder("jpeg", &Encoder{})
}
