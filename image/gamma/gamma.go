// Package gamma provides gamma imageserver/image.Processor implementations.
package gamma

import (
	"context"
	"image"
	"image/draw"
	"math"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	imageserver_image_internal "github.com/pierrre/imageserver/image/internal"
	"github.com/pierrre/imageutil"
)

// Processor is a imageserver/image.Processor implementation that applies gamma transformation.
type Processor struct {
	vals        [1 << 16]uint16
	newDrawable func(image.Image) draw.Image
}

// NewProcessor creates a Processor.
//
// "highQuality" indicates if the Processor return a NRGBA64 Image or an Image with the same quality as the given Image.
func NewProcessor(gamma float64, highQuality bool) *Processor {
	prc := new(Processor)
	gammaInv := 1 / gamma
	for i := range prc.vals {
		prc.vals[i] = uint16(math.Pow(float64(i)/65535, gammaInv)*65535 + 0.5)
	}
	if highQuality {
		prc.newDrawable = func(p image.Image) draw.Image {
			return image.NewNRGBA64(p.Bounds())
		}
	} else {
		prc.newDrawable = imageserver_image_internal.NewDrawable
	}
	return prc
}

// Process implements imageserver/image.Processor.
//
// It doesn't return an error.
func (prc *Processor) Process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	out := prc.newDrawable(nim)
	bd := nim.Bounds().Intersect(out.Bounds())
	at := imageutil.NewAtFunc(nim)
	set := imageutil.NewSetFunc(out)
	imageutil.Parallel1D(ctx, bd, func(ctx context.Context, bd image.Rectangle) {
		for y := bd.Min.Y; y < bd.Max.Y; y++ {
			for x := bd.Min.X; x < bd.Max.X; x++ {
				r, g, b, a := at(x, y)
				r, g, b, a = imageutil.RGBAToNRGBA(r, g, b, a)
				r = uint32(prc.vals[uint16(r)])
				g = uint32(prc.vals[uint16(g)])
				b = uint32(prc.vals[uint16(b)])
				r, g, b, a = imageutil.NRGBAToRGBA(r, g, b, a)
				set(x, y, r, g, b, a)
			}
		}
	})
	return out, nil
}

// Change implements imageserver/image.Processor.
func (prc *Processor) Change(params imageserver.Params) bool {
	return true
}

const correct = 2.2

// CorrectionProcessor is a imageserver/image.Processor implementation that corrects gamma for a sub Processor.
//
// Steps:
//  - apply gamma of 1/2.2 (darken)
//  - call the sub Processor
//  - apply gamma of 2.2 (lighten)
//
// Internally, it uses NRGBA64 high quality Image, to avoid loss of information.
//
// The CorrectionProcessor can be enabled/disabled with the "gamma_correction" (bool) param.
type CorrectionProcessor struct {
	imageserver_image.Processor
	enabled bool
	before  *Processor
	after   *Processor
}

// NewCorrectionProcessor creates a CorrectionProcessor.
//
// "enabled" indicated if the CorrectionProcessor is enabled by default.
func NewCorrectionProcessor(prc imageserver_image.Processor, enabled bool) *CorrectionProcessor {
	return &CorrectionProcessor{
		Processor: prc,
		enabled:   enabled,
		before:    NewProcessor(1/correct, true),
		after:     NewProcessor(correct, true),
	}
}

// Process implements imageserver/image.Processor.
func (prc *CorrectionProcessor) Process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	enabled, err := prc.isEnabled(params)
	if err != nil {
		return nil, err
	}
	if enabled {
		return prc.process(ctx, nim, params)
	}
	return prc.Processor.Process(ctx, nim, params)
}

func (prc *CorrectionProcessor) isEnabled(params imageserver.Params) (bool, error) {
	if params.Has("gamma_correction") {
		return params.GetBool("gamma_correction")
	}
	return prc.enabled, nil
}

func (prc *CorrectionProcessor) process(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
	original := nim
	nim, _ = prc.before.Process(ctx, nim, params)
	nim, err := prc.Processor.Process(ctx, nim, params)
	if err != nil {
		return nil, err
	}
	nim, _ = prc.after.Process(ctx, nim, params)
	if isHighQuality(nim) && !isHighQuality(original) {
		newNim := imageserver_image_internal.NewDrawableSize(original, nim.Bounds())
		imageserver_image_internal.Copy(ctx, newNim, nim)
		nim = newNim
	}
	return nim, nil
}

func isHighQuality(p image.Image) bool {
	switch p.(type) {
	case *image.RGBA64, *image.NRGBA64:
		return true
	default:
		return false
	}
}
