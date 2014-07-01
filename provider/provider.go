package provider

import (
	"github.com/pierrre/imageserver"
)

// Provider is an Image provider
type Provider interface {
	Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error)
}

// ProviderFunc is a Provider func
type ProviderFunc func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error)

// Get call the func
func (f ProviderFunc) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return f(source, parameters)
}
