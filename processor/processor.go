package processor

import (
	"github.com/pierrre/imageserver"
)

// Processor represents an Image processor
type Processor interface {
	Process(*imageserver.Image, imageserver.Parameters) (*imageserver.Image, error)
}

// ProcessorFunc is a Processor func
type ProcessorFunc func(*imageserver.Image, imageserver.Parameters) (*imageserver.Image, error)

// Process calls the func
func (f ProcessorFunc) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return f(image, parameters)
}
