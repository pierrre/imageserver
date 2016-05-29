package image

import (
	"image"

	"github.com/pierrre/imageserver"
)

// Provider returns a Go Image.
type Provider interface {
	Get(imageserver.Params) (image.Image, error)
}

// ProviderFunc is a Provider func.
type ProviderFunc func(imageserver.Params) (image.Image, error)

// Get implements Provider.
func (f ProviderFunc) Get(params imageserver.Params) (image.Image, error) {
	return f(params)
}

// ProcessorProvider is a Provider implementation that processes the Image.
type ProcessorProvider struct {
	Provider
	Processor Processor
}

// Get implements Provider.
func (prv *ProcessorProvider) Get(params imageserver.Params) (image.Image, error) {
	nim, err := prv.Provider.Get(params)
	if err != nil {
		return nil, err
	}
	nim, err = prv.Processor.Process(nim, params)
	if err != nil {
		return nil, err
	}
	return nim, nil
}
