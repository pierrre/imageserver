// Package limit provides an Image Processor that limits the number of concurrent executions
package limit

import (
	"github.com/pierrre/imageserver"
)

// LimitProcessor represents an ImageProcessor that limits the number of concurrent executions
//
// It wraps an Image Processor and use a buffered channel to limit the number of concurrent executions.
type LimitProcessor struct {
	limitChan chan bool
	processor imageserver.Processor
}

// New creates a LimitProcessor
func New(limit uint, processor imageserver.Processor) imageserver.Processor {
	return &LimitProcessor{
		limitChan: make(chan bool, limit),
		processor: processor,
	}
}

// Process forwards the call to the wrapped Image Processor and limits the number of concurrent executions
func (processor *LimitProcessor) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	processor.limitChan <- true
	defer func() {
		<-processor.limitChan
	}()
	return processor.processor.Process(image, parameters)
}
