// Package crop provides a imageserver/image.Processor implementation that allows to crop Image.
package crop

import (
	"context"
	"fmt"
	"image"

	"github.com/pierrre/imageserver"
)

const param = "crop"

// Processor is a imageserver/image.Processor implementation that allows to crop Image.
//
// All params are extracted from the "crop" node param and are mandatory:
//  - min_x: top-left X
//  - min_y: top-left Y
//  - max_x: bottom-right X
//  - max_y: bottom-right Y
type Processor struct{}

// Process implements imageserver/image.Processor.
func (prc *Processor) Process(ctx context.Context, im image.Image, params imageserver.Params) (image.Image, error) {
	if !params.Has(param) {
		return im, nil
	}
	params, err := params.GetParams(param)
	if err != nil {
		return nil, err
	}
	im, err = prc.process(im, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = param + "." + err.Param
		}
		return nil, err
	}
	return im, nil
}

func (prc *Processor) process(im image.Image, params imageserver.Params) (image.Image, error) {
	bds, err := prc.getBounds(params)
	if err != nil {
		return nil, err
	}
	return prc.crop(im, bds)
}

func (prc *Processor) getBounds(params imageserver.Params) (image.Rectangle, error) {
	var bds image.Rectangle
	var err error
	bds.Min.X, err = params.GetInt("min_x")
	if err != nil {
		return image.ZR, err
	}
	bds.Min.Y, err = params.GetInt("min_y")
	if err != nil {
		return image.ZR, err
	}
	bds.Max.X, err = params.GetInt("max_x")
	if err != nil {
		return image.ZR, err
	}
	bds.Max.Y, err = params.GetInt("max_y")
	if err != nil {
		return image.ZR, err
	}
	return bds, nil
}

func (prc *Processor) crop(im image.Image, bds image.Rectangle) (image.Image, error) {
	type SubImage interface {
		image.Image
		SubImage(image.Rectangle) image.Image
	}
	im2, ok := im.(SubImage)
	if !ok {
		return nil, &imageserver.ImageError{
			Message: fmt.Sprintf("crop: image type %T not supported: method SubImage not found", im),
		}
	}
	return im2.SubImage(bds), nil
}

// Change implements imageserver/image.Processor.
func (prc *Processor) Change(params imageserver.Params) bool {
	return params.Has(param)
}
