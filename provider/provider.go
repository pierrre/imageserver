package provider

import (
	"github.com/pierrre/imageserver"
)

// Provider is an Image provider
type Provider interface {
	Get(source interface{}, params imageserver.Params) (*imageserver.Image, error)
}

// Func is a Provider func
type Func func(source interface{}, params imageserver.Params) (*imageserver.Image, error)

// Get call the func
func (f Func) Get(source interface{}, params imageserver.Params) (*imageserver.Image, error) {
	return f(source, params)
}
