package provider

import (
	"github.com/pierrre/imageserver"
)

// ProviderImageServer is a Provider ImageServer
type ProviderImageServer struct {
	Provider Provider
}

// Get get an Image from the Provider using the "source" parameter
func (pis *ProviderImageServer) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	source, err := parameters.Get("source")
	if err != nil {
		return nil, imageserver.NewError("Missing source parameter")
	}

	image, err := pis.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	return image, nil
}
