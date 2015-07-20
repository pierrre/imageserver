package image

import (
	"image"

	"github.com/pierrre/imageserver"
)

// Provider represents a Go Image provider.
type Provider interface {
	Get(params imageserver.Params) (image.Image, error)
}

// ProviderFunc is a Provider func.
type ProviderFunc func(imageserver.Params) (image.Image, error)

// Get implements Provider
func (f ProviderFunc) Get(params imageserver.Params) (image.Image, error) {
	return f(params)
}

// StaticProvider is a Go Image Provider that always return the same result.
type StaticProvider struct {
	Image image.Image
	Error error
}

// Get implements Provider
func (prv *StaticProvider) Get(params imageserver.Params) (image.Image, error) {
	return prv.Image, prv.Error
}

// ProcessorProvider is a Go Image provider that processes the Image.
type ProcessorProvider struct {
	Provider
	Processor Processor
}

// Get implements Provider
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

// ProviderServer is an Image Server for an Go Image Provider.
type ProviderServer struct {
	Provider Provider
}

// Get implements Server
func (srv *ProviderServer) Get(params imageserver.Params) (*imageserver.Image, error) {
	enc, format, err := getEncoderFormat("", params)
	if err != nil {
		return nil, err
	}
	nim, err := srv.Provider.Get(params)
	if err != nil {
		return nil, err
	}
	im, err := encode(nim, format, enc, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}
