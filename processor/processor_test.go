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
	ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return image, nil
	})(testdata.Small, make(imageserver.Parameters))
}
