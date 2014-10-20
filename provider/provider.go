package provider

import (
	"github.com/pierrre/imageserver"
)

// Provider is an Image provider
type Provider interface {
	Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error)
}

// Func is a Provider func
type Func func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error)

// Get call the func
func (f Func) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return f(source, parameters)
}
