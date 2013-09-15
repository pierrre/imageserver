// Limit processor
package limit

import (
	"github.com/pierrre/imageserver"
)

// Limit concurrent usage of a processor
type LimitProcessor struct {
	limitChan chan bool
	processor imageserver.Processor
}

func New(limit uint, processor imageserver.Processor) imageserver.Processor {
	return &LimitProcessor{
		limitChan: make(chan bool, limit),
		processor: processor,
	}
}

func (processor *LimitProcessor) Process(inImage *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	processor.limitChan <- true
	defer func() {
		<-processor.limitChan
	}()
	return processor.processor.Process(inImage, parameters)
}
