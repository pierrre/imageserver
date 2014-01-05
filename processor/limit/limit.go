// Package limit provides an Image Processor that limits the number of concurrent executions
package limit

import (
	"github.com/pierrre/imageserver"
)

// LimitProcessor represents an ImageProcessor that limits the number of concurrent executions
//
// It wraps an Image Processor and use a buffered channel to limit the number of concurrent executions.
type LimitProcessor struct {
	imageserver.Processor
	LimitChan chan bool
}

// New creates a LimitProcessor
func New(processor imageserver.Processor, limit uint) imageserver.Processor {
	return &LimitProcessor{
		Processor: processor,
		LimitChan: make(chan bool, limit),
	}
}

// Process forwards the call to the wrapped Image Processor and limits the number of concurrent executions
func (processor *LimitProcessor) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	processor.LimitChan <- true
	defer func() {
		<-processor.LimitChan
	}()
	return processor.Processor.Process(image, parameters)
}
