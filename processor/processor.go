package processor

import (
	"github.com/pierrre/imageserver"
)

// Processor represents an Image processor
type Processor interface {
	Process(*imageserver.Image, imageserver.Params) (*imageserver.Image, error)
}

// Func is a Processor func
type Func func(*imageserver.Image, imageserver.Params) (*imageserver.Image, error)

// Process calls the func
func (f Func) Process(image *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	return f(image, params)
}

// List represents a list of Image Processor
type List []Processor

// Process processes the Image with the list of Image Processor
func (l List) Process(image *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	var err error
	for _, p := range l {
		image, err = p.Process(image, params)
		if err != nil {
			return nil, err
		}
	}
	return image, nil
}

/*
NewLimit creates a new Processor that limits the number of concurrent executions.

It uses a buffered channel to limit the number of concurrent executions.
*/
func NewLimit(processor Processor, limit uint) Processor {
	limitCh := make(chan struct{}, limit)
	return Func(func(image *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
		limitCh <- struct{}{}
		defer func() {
			<-limitCh
		}()
		return processor.Process(image, params)
	})
}
