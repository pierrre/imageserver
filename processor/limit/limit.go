// Package limit provides an Image Processor that limits the number of concurrent executions
package limit

import (
	"github.com/pierrre/imageserver"
	imageserver_processor "github.com/pierrre/imageserver/processor"
)

// LimitProcessor represents an Image Processor that limits the number of concurrent executions
//
// It wraps an Image Processor and use a buffered channel to limit the number of concurrent executions.
type LimitProcessor struct {
	Processor imageserver_processor.Processor
	LimitChan chan struct{}
}

// New creates a LimitProcessor
func New(processor imageserver_processor.Processor, limit uint) *LimitProcessor {
	return &LimitProcessor{
		Processor: processor,
		LimitChan: make(chan struct{}, limit),
	}
}

// Process forwards the call to the wrapped Image Processor and limits the number of concurrent executions
func (processor *LimitProcessor) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	processor.LimitChan <- struct{}{}
	defer func() {
		<-processor.LimitChan
	}()
	return processor.Processor.Process(image, parameters)
}
