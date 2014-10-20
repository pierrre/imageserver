package processor

import (
	"github.com/pierrre/imageserver"
)

// Processor represents an Image processor
type Processor interface {
	Process(*imageserver.Image, imageserver.Parameters) (*imageserver.Image, error)
}

// Func is a Processor func
type Func func(*imageserver.Image, imageserver.Parameters) (*imageserver.Image, error)

// Process calls the func
func (f Func) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return f(image, parameters)
}

// List represents a list of Image Processor
type List []Processor

// Process processes the Image with the list of Image Processor
func (l List) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	var err error
	for _, p := range l {
		image, err = p.Process(image, parameters)
		if err != nil {
			return nil, err
		}
	}
	return image, nil
}

/*
Limit represents an Image Processor that limits the number of concurrent executions.

It wraps an Image Processor and use a buffered channel to limit the number of concurrent executions.
*/
type Limit struct {
	Processor
	limitChan chan struct{}
}

// NewLimit creates a Limit
func NewLimit(processor Processor, limit uint) *Limit {
	return &Limit{
		Processor: processor,
		limitChan: make(chan struct{}, limit),
	}
}

// Process forwards the call to the wrapped Image Processor and limits the number of concurrent executions
func (processor *Limit) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	processor.limitChan <- struct{}{}
	defer func() {
		<-processor.limitChan
	}()
	return processor.Processor.Process(image, parameters)
}
