package processor

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestProcessorFuncInterface(t *testing.T) {
	var _ Processor = ProcessorFunc(nil)
}

func TestProcessorFunc(t *testing.T) {
	pf := ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return image, nil
	})
	pf.Process(testdata.Small, make(imageserver.Parameters))
}
