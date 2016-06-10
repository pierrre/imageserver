package image

import (
	"context"

	"github.com/pierrre/imageserver"
)

// Handler is a imageserver.Handler implementation that uses Go "image" package.
//
// It supports format conversion and processing.
// It uses the "format" param to determine which Encoder is used.
//
// If there is nothing to do, Handler does not decode the Image or call the Processor.
type Handler struct {
	Processor Processor // Optional Processor
}

// Handle implements imageserver.Handler.
func (hdr *Handler) Handle(ctx context.Context, im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	enc, format, err := getEncoderFormat(im.Format, params)
	if err != nil {
		if _, ok := err.(*imageserver.ParamError); !ok {
			err = &imageserver.ImageError{Message: err.Error()}
		}
		return nil, err
	}
	if !hdr.change(im, format, enc, params) {
		return im, nil
	}
	nim, err := Decode(im)
	if err != nil {
		return nil, err
	}
	if hdr.Processor != nil {
		nim, err = hdr.Processor.Process(ctx, nim, params)
		if err != nil {
			return nil, err
		}
	}
	im, err = encode(nim, format, enc, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (hdr *Handler) change(im *imageserver.Image, format string, enc Encoder, params imageserver.Params) bool {
	if format != im.Format {
		return true
	}
	if hdr.Processor != nil && hdr.Processor.Change(params) {
		return true
	}
	if enc.Change(params) {
		return true
	}
	return false
}
