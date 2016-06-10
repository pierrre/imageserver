// Package nfntresize provides a nfnt/resize imageserver/image.Processor implementation.
package nfntresize

import (
	"context"
	"fmt"
	"image"

	"github.com/nfnt/resize"
	"github.com/pierrre/imageserver"
)

const (
	param = "nfntresize"
)

// Processor is a nfnt/resize imageserver/image.Processor implementation.
//
// All params are extracted from the "graphicsmagick" node param and are optionals:
//  - width
//  - height
//  - mode: resize mode
//      possible values:
//      - resize (default): see github.com/nfnt/resize.Resize
//      - thumbnail: see github.com/nfnt/resize.Thumbnail
//  - interpolation: interpolation method
//      possible values:
//      - nearest_neighbor (default)
//      - bilinear
//      - bicubic
//      - mitchell_netravali
//      - lanczos2
//      - lanczos3
type Processor struct {
	DefaultInterpolation resize.InterpolationFunction
	MaxWidth             int
	MaxHeight            int
}

// Process implements imageserver/image.Processor.
func (prc *Processor) Process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	if !params.Has(param) {
		return nim, nil
	}
	params, err := params.GetParams(param)
	if err != nil {
		return nil, err
	}
	if params.Empty() {
		return nim, nil
	}
	nim, err = prc.process(nim, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = fmt.Sprintf("%s.%s", param, err.Param)
		}
		return nil, err
	}
	return nim, nil
}

func (prc *Processor) process(nim image.Image, params imageserver.Params) (image.Image, error) {
	width, height, err := prc.getSize(params)
	if err != nil {
		return nil, err
	}
	if width == 0 && height == 0 {
		return nim, nil
	}
	interp, err := prc.getInterpolation(params)
	if err != nil {
		return nil, err
	}
	mode, err := getModeFunc(params)
	if err != nil {
		return nil, err
	}
	nim = mode(width, height, nim, interp)
	return nim, nil
}

func (prc *Processor) getSize(params imageserver.Params) (uint, uint, error) {
	w, err := getDimension("width", prc.MaxWidth, params)
	if err != nil {
		return 0, 0, err
	}
	h, err := getDimension("height", prc.MaxHeight, params)
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}

func getDimension(name string, max int, params imageserver.Params) (uint, error) {
	if !params.Has(name) {
		return 0, nil
	}
	d, err := params.GetInt(name)
	if err != nil {
		return 0, err
	}
	if d < 0 {
		return 0, &imageserver.ParamError{Param: name, Message: "must be greater than or equal to 0"}
	}
	if max > 0 && d > max {
		return 0, &imageserver.ParamError{Param: name, Message: fmt.Sprintf("must be less than or equal to %d", max)}
	}
	return uint(d), nil
}

func (prc *Processor) getInterpolation(params imageserver.Params) (resize.InterpolationFunction, error) {
	if !params.Has("interpolation") {
		return prc.DefaultInterpolation, nil
	}
	interpolation, err := params.GetString("interpolation")
	if err != nil {
		return 0, err
	}
	switch interpolation {
	case "nearest_neighbor":
		return resize.NearestNeighbor, nil
	case "bilinear":
		return resize.Bilinear, nil
	case "bicubic":
		return resize.Bicubic, nil
	case "mitchell_netravali":
		return resize.MitchellNetravali, nil
	case "lanczos2":
		return resize.Lanczos2, nil
	case "lanczos3":
		return resize.Lanczos3, nil
	default:
		return 0, &imageserver.ParamError{Param: "interpolation", Message: "invalid value"}
	}
}

type modeFunc func(uint, uint, image.Image, resize.InterpolationFunction) image.Image

func getModeFunc(params imageserver.Params) (modeFunc, error) {
	if !params.Has("mode") {
		return resize.Resize, nil
	}
	mode, err := params.GetString("mode")
	if err != nil {
		return nil, err
	}
	switch mode {
	case "resize":
		return resize.Resize, nil
	case "thumbnail":
		return resize.Thumbnail, nil
	default:
		return nil, &imageserver.ParamError{Param: "mode", Message: "invalid value"}
	}
}

// Change implements imageserver/image.Processor.
func (prc *Processor) Change(params imageserver.Params) bool {
	if !params.Has(param) {
		return false
	}
	params, err := params.GetParams(param)
	if err != nil {
		return true
	}
	if params.Empty() {
		return false
	}
	if params.Has("width") {
		return true
	}
	if params.Has("height") {
		return true
	}
	return false
}
