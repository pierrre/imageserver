// Package imageserver provides an image server
package imageserver

// ImageServer is an interface for an Image server
type ImageServer interface {
	Get(Parameters) (*Image, error)
}

// ImageServerFunc is a ImageServer func
type ImageServerFunc func(parameters Parameters) (*Image, error)

// Get calls the func
func (f ImageServerFunc) Get(parameters Parameters) (*Image, error) {
	return f(parameters)
}
