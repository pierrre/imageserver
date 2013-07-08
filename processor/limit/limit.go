package limit

import (
	"github.com/pierrre/imageserver"
)

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

func (processor *LimitProcessor) Process(inImage *imageserver.Image, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	processor.limitChan <- true
	defer func() {
		<-processor.limitChan
	}()
	return processor.processor.Process(inImage, parameters)
}
