package gift

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/gift"
	"github.com/pierrre/imageserver"
	imageserver_image_internal "github.com/pierrre/imageserver/image/internal"
)

const (
	rotateParam = "gift_rotate"
)

// RotateProcessor is a imageserver/image.Processor implementation that rotates the Image with GIFT.
type RotateProcessor struct {
	DefaultInterpolation gift.Interpolation
}

// Process implements imageserver/image.Processor.
func (prc *RotateProcessor) Process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	if !params.Has(rotateParam) {
		return nim, nil
	}
	params, err := params.GetParams(rotateParam)
	if err != nil {
		return nil, err
	}
	if params.Empty() {
		return nim, nil
	}
	nim, err = prc.process(nim, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = fmt.Sprintf("%s.%s", rotateParam, err.Param)
		}
		return nil, err
	}
	return nim, nil
}

func (prc *RotateProcessor) process(nim image.Image, params imageserver.Params) (image.Image, error) {
	rot, err := prc.getRotation(params)
	if err != nil {
		return nil, err
	}
	if rot == 0 {
		return nim, nil
	}
	f, err := prc.getFilter(rot, params)
	if err != nil {
		return nil, err
	}
	g := gift.New(f)
	out := imageserver_image_internal.NewDrawableSize(nim, g.Bounds(nim.Bounds()))
	g.Draw(out, nim)
	return out, nil
}

func (prc *RotateProcessor) getRotation(params imageserver.Params) (float32, error) {
	if !params.Has("rotation") {
		return 0, nil
	}
	rot, err := params.GetFloat("rotation")
	if err != nil {
		return 0, err
	}
	if rot < 0 {
		rot = math.Mod(rot, 360) + 360
	}
	if rot >= 360 {
		rot = math.Mod(rot, 360)
	}
	return float32(rot), nil
}

func (prc *RotateProcessor) getFilter(rot float32, params imageserver.Params) (gift.Filter, error) {
	switch rot {
	case 90:
		return gift.Rotate90(), nil
	case 180:
		return gift.Rotate180(), nil
	case 270:
		return gift.Rotate270(), nil
	}
	bkg, err := prc.getBackground(params)
	if err != nil {
		return nil, err
	}
	interp, err := prc.getInterpolation(params)
	if err != nil {
		return nil, err
	}
	return gift.Rotate(rot, bkg, interp), nil
}

func (prc *RotateProcessor) getBackground(params imageserver.Params) (color.Color, error) {
	if !params.Has("background") {
		return color.Transparent, nil
	}
	s, err := params.GetString("background")
	if err != nil {
		return nil, err
	}
	c, err := parseHexColor(s)
	if err != nil {
		return nil, &imageserver.ParamError{Param: "background", Message: err.Error()}
	}
	return c, nil
}

func (prc *RotateProcessor) getInterpolation(params imageserver.Params) (gift.Interpolation, error) {
	if !params.Has("interpolation") {
		return prc.DefaultInterpolation, nil
	}
	interp, err := params.GetString("interpolation")
	if err != nil {
		return 0, err
	}
	switch interp {
	case "nearest_neighbor":
		return gift.NearestNeighborInterpolation, nil
	case "linear":
		return gift.LinearInterpolation, nil
	case "cubic":
		return gift.CubicInterpolation, nil
	}
	return 0, &imageserver.ParamError{Param: "interpolation", Message: "invalid value"}
}

// Change implements imageserver/image.Processor.
func (prc *RotateProcessor) Change(params imageserver.Params) bool {
	if !params.Has(rotateParam) {
		return false
	}
	params, err := params.GetParams(rotateParam)
	if err != nil {
		return true
	}
	if params.Empty() {
		return false
	}
	if params.Has("rotation") {
		return true
	}
	return false
}

func parseHexColor(s string) (color.Color, error) {
	if len(s) > 8 {
		return nil, fmt.Errorf("too long: %d", len(s))
	}
	is, err := hexStringToInts(s)
	if err != nil {
		return nil, err
	}
	switch len(is) {
	case 3:
		return color.NRGBA{
			R: is[0] * 0x11,
			G: is[1] * 0x11,
			B: is[2] * 0x11,
			A: 0xff,
		}, nil
	case 4:
		return color.NRGBA{
			A: is[0] * 0x11,
			R: is[1] * 0x11,
			G: is[2] * 0x11,
			B: is[3] * 0x11,
		}, nil
	case 6:
		return color.NRGBA{
			R: is[0]*0x10 + is[1],
			G: is[2]*0x10 + is[3],
			B: is[4]*0x10 + is[5],
			A: 0xff,
		}, nil
	case 8:
		return color.NRGBA{
			A: is[0]*0x10 + is[1],
			R: is[2]*0x10 + is[3],
			G: is[4]*0x10 + is[5],
			B: is[6]*0x10 + is[7],
		}, nil
	default:
		return nil, fmt.Errorf("invalid length: %d", len(s))
	}
}

func hexStringToInts(s string) ([]uint8, error) {
	var res []uint8
	for i, r := range s {
		var v uint8
		if r >= '0' && r <= '9' {
			v = uint8(0x0 + r - '0')
		} else if r >= 'a' && r <= 'f' {
			v = uint8(0xa + r - 'a')
		} else if r >= 'A' && r <= 'F' {
			v = uint8(0xa + r - 'A')
		} else {
			return nil, fmt.Errorf("invalid character '%c' at position %d", r, i)
		}
		res = append(res, v)
	}
	return res, nil
}
