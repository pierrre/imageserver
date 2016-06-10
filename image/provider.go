package image

import (
	"context"
	"image"

	"github.com/pierrre/imageserver"
)

// Provider returns a Go Image.
type Provider interface {
	Get(context.Context, imageserver.Params) (image.Image, error)
}

// ProviderFunc is a Provider func.
type ProviderFunc func(context.Context, imageserver.Params) (image.Image, error)

// Get implements Provider.
func (f ProviderFunc) Get(ctx context.Context, params imageserver.Params) (image.Image, error) {
	return f(ctx, params)
}

// ProcessorProvider is a Provider implementation that processes the Image.
type ProcessorProvider struct {
	Provider
	Processor Processor
}

// Get implements Provider.
func (prv *ProcessorProvider) Get(ctx context.Context, params imageserver.Params) (image.Image, error) {
	nim, err := prv.Provider.Get(ctx, params)
	if err != nil {
		return nil, err
	}
	nim, err = prv.Processor.Process(ctx, nim, params)
	if err != nil {
		return nil, err
	}
	return nim, nil
}
