// Package imageserver provides an Image server
package imageserver

// Server is an interface for an Image server
type Server interface {
	Get(Params) (*Image, error)
}

// ServerFunc is a Server func
type ServerFunc func(params Params) (*Image, error)

// Get calls the func
func (f ServerFunc) Get(params Params) (*Image, error) {
	return f(params)
}
