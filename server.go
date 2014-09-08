// Package imageserver provides an Image server
package imageserver

// Server is an interface for an Image server
type Server interface {
	Get(Parameters) (*Image, error)
}

// ServerFunc is a Server func
type ServerFunc func(parameters Parameters) (*Image, error)

// Get calls the func
func (f ServerFunc) Get(parameters Parameters) (*Image, error) {
	return f(parameters)
}
