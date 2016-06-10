package gift

import (
	"context"
	"fmt"
	"image"

	"github.com/disintegration/gift"
	"github.com/pierrre/imageserver"
	imageserver_image_internal "github.com/pierrre/imageserver/image/internal"
)

const (
	resizeParam = "gift_resize"
)

// ResizeProcessor is a imageserver/image.Processor implementation that resizes the Image with GIFT.
//
// All params are extracted from the "gift_resize" node param and are optionals:
//  - width
//  - height
//  - mode: resize mode
//      possible values:
//      - <no value> (default): see github.com/disintegration/gift.Resize
//      - fit: see github.com/disintegration/gift.ResizeToFit
//      - fill: see github.com/disintegration/gift.ResizeToFill
//  - resampling: resampling method
//      possible values:
//      - nearest_neighbor (default)
//      - box
//      - linear
//      - cubic
//      - lanczos
type ResizeProcessor struct {
	DefaultResampling gift.Resampling
	MaxWidth          int
	MaxHeight         int
}

// Process implements imageserver/image.Processor.
func (prc *ResizeProcessor) Process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	if !params.Has(resizeParam) {
		return nim, nil
	}
	params, err := params.GetParams(resizeParam)
	if err != nil {
		return nil, err
	}
	if params.Empty() {
		return nim, nil
	}
	nim, err = prc.process(nim, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = fmt.Sprintf("%s.%s", resizeParam, err.Param)
		}
		return nil, err
	}
	return nim, nil
}

func (prc *ResizeProcessor) process(nim image.Image, params imageserver.Params) (image.Image, error) {
	width, height, err := prc.getSize(params)
	if err != nil {
		return nil, err
	}
	if width == 0 && height == 0 {
		return nim, err
	}
	f, err := prc.getFilter(width, height, params)
	if err != nil {
		return nil, err
	}
	g := gift.New(f)
	out := imageserver_image_internal.NewDrawableSize(nim, g.Bounds(nim.Bounds()))
	g.Draw(out, nim)
	return out, nil
}

func (prc *ResizeProcessor) getSize(params imageserver.Params) (int, int, error) {
	w, err := prc.getDimension("width", prc.MaxWidth, params)
	if err != nil {
		return 0, 0, err
	}
	h, err := prc.getDimension("height", prc.MaxHeight, params)
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}

func (prc *ResizeProcessor) getDimension(name string, max int, params imageserver.Params) (int, error) {
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
	return d, nil
}

func (prc *ResizeProcessor) getFilter(width, height int, params imageserver.Params) (gift.Filter, error) {
	rsp, err := prc.getResampling(params)
	if err != nil {
		return nil, err
	}
	if !params.Has("mode") || width == 0 || height == 0 {
		return gift.Resize(width, height, rsp), nil
	}
	mode, err := params.GetString("mode")
	if err != nil {
		return nil, err
	}
	switch mode {
	case "fit":
		return gift.ResizeToFit(width, height, rsp), nil
	case "fill":
		return gift.ResizeToFill(width, height, rsp, gift.CenterAnchor), nil
	}
	return nil, &imageserver.ParamError{Param: "mode", Message: "invalid value"}
}

func (prc *ResizeProcessor) getResampling(params imageserver.Params) (gift.Resampling, error) {
	if !params.Has("resampling") {
		if prc.DefaultResampling != nil {
			return prc.DefaultResampling, nil
		}
		return gift.NearestNeighborResampling, nil
	}
	rsp, err := params.GetString("resampling")
	if err != nil {
		return nil, err
	}
	switch rsp {
	case "nearest_neighbor":
		return gift.NearestNeighborResampling, nil
	case "box":
		return gift.BoxResampling, nil
	case "linear":
		return gift.LinearResampling, nil
	case "cubic":
		return gift.CubicResampling, nil
	case "lanczos":
		return gift.LanczosResampling, nil
	}
	return nil, &imageserver.ParamError{Param: "resampling", Message: "invalid value"}
}

// Change implements imageserver/image.Processor.
func (prc *ResizeProcessor) Change(params imageserver.Params) bool {
	if !params.Has(resizeParam) {
		return false
	}
	params, err := params.GetParams(resizeParam)
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
