package gif

import (
	"bytes"
	"context"
	"fmt"
	"image/gif"

	"github.com/pierrre/imageserver"
)

// Handler is a GIF imageserver.Handler implementation.
//
// Steps:
//  - decode the GIF image
//  - processes the GIF image with a Processor from this package
//  - encode the image to GIF
//
// If there is nothing to do, Handler does not decode the GIF image or call the Processor.
type Handler struct {
	Processor Processor
}

// Handle implements imageserver.Handler.
func (hdr *Handler) Handle(ctx context.Context, im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	if im.Format != "gif" {
		return nil, &imageserver.ImageError{Message: fmt.Sprintf("image format is not gif: %s", im.Format)}
	}
	if !hdr.Processor.Change(params) {
		return im, nil
	}
	g, err := gif.DecodeAll(bytes.NewReader(im.Data))
	if err != nil {
		return nil, &imageserver.ImageError{Message: err.Error()}
	}
	g, err = hdr.Processor.Process(ctx, g, params)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = gif.EncodeAll(buf, g)
	if err != nil {
		return nil, &imageserver.ImageError{Message: err.Error()}
	}
	im = &imageserver.Image{
		Format: "gif",
		Data:   buf.Bytes(),
	}
	return im, nil
}

// FallbackHandler is a imageserver.Handler implementation that allows to switch between a Handler of this package, or a fallback Handler.
//
// If the Image format and the "format" param are equal to "gif", the Handler of this package is used.
// Otherwise, the fallback Handler is used.
type FallbackHandler struct {
	*Handler
	Fallback imageserver.Handler
}

// Handle implements imageserver.Handler.
func (hdr *FallbackHandler) Handle(ctx context.Context, im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	h, err := hdr.getHandler(im, params)
	if err != nil {
		return nil, err
	}
	return h.Handle(ctx, im, params)
}

func (hdr *FallbackHandler) getHandler(im *imageserver.Image, params imageserver.Params) (imageserver.Handler, error) {
	if im.Format != "gif" {
		return hdr.Fallback, nil
	}
	if !params.Has("format") {
		return hdr.Handler, nil
	}
	format, err := params.GetString("format")
	if err != nil {
		return nil, err
	}
	if format != "gif" {
		return hdr.Fallback, nil
	}
	return hdr.Handler, nil
}
