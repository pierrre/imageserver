package gif

import (
	"context"
	"image"
	"image/gif"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	imageserver_image_internal "github.com/pierrre/imageserver/image/internal"
)

// Processor processes a GIF image.
type Processor interface {
	Process(context.Context, *gif.GIF, imageserver.Params) (*gif.GIF, error)
	imageserver_image.Changer
}

// SimpleProcessor is a Processor implementation that processes each frames with the sub imageserver/image.Processor.
type SimpleProcessor struct {
	imageserver_image.Processor
}

// Process implements Processor.
func (prc *SimpleProcessor) Process(ctx context.Context, g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
	out := new(gif.GIF)
	var err error
	out.Image, err = prc.processImages(ctx, g.Image, params)
	if err != nil {
		return nil, err
	}
	out.Delay = make([]int, len(g.Delay))
	copy(out.Delay, g.Delay)
	out.LoopCount = g.LoopCount
	if g.Disposal != nil {
		out.Disposal = make([]byte, len(g.Disposal))
		copy(out.Disposal, g.Disposal)
	}
	out.Config.ColorModel = g.Config.ColorModel
	for _, p := range out.Image {
		if p.Rect.Max.X > out.Config.Width {
			out.Config.Width = p.Rect.Max.X
		}
		if p.Rect.Max.Y > out.Config.Height {
			out.Config.Height = p.Rect.Max.Y
		}
	}
	out.BackgroundIndex = g.BackgroundIndex
	return out, nil
}

func (prc *SimpleProcessor) processImages(ctx context.Context, ps []*image.Paletted, params imageserver.Params) ([]*image.Paletted, error) {
	out := make([]*image.Paletted, len(ps))
	for i, p := range ps {
		var err error
		out[i], err = prc.processImage(ctx, p, params)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (prc *SimpleProcessor) processImage(ctx context.Context, p *image.Paletted, params imageserver.Params) (*image.Paletted, error) {
	tmp, err := prc.Processor.Process(ctx, p, params)
	if err != nil {
		return nil, err
	}
	out, ok := tmp.(*image.Paletted)
	if !ok {
		out = image.NewPaletted(tmp.Bounds(), p.Palette)
		imageserver_image_internal.Copy(ctx, out, tmp)
	}
	return out, nil
}

// ProcessorFunc is a Processor func.
type ProcessorFunc func(context.Context, *gif.GIF, imageserver.Params) (*gif.GIF, error)

// Process implements Processor.
func (f ProcessorFunc) Process(ctx context.Context, g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
	return f(ctx, g, params)
}

// Change implements Processor.
func (f ProcessorFunc) Change(params imageserver.Params) bool {
	return true
}
