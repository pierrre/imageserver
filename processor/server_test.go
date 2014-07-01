package processor

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestProcessorImageServerInterface(t *testing.T) {
	var _ imageserver.ImageServer = &ProcessorImageServer{}
}

func TestProcessorImageServer(t *testing.T) {
	pis := &ProcessorImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Small, nil
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return image, nil
		}),
	}
	image, err := pis.Get(make(imageserver.Parameters))
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("no image")
	}
}

func TestProcessorImageServerErrorImageServer(t *testing.T) {
	pis := &ProcessorImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return image, nil
		}),
	}
	_, err := pis.Get(make(imageserver.Parameters))
	if err == nil {
		t.Fatal("no error")
	}
}

func TestProcessorImageServerErrorProcessor(t *testing.T) {
	pis := &ProcessorImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Small, nil
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := pis.Get(make(imageserver.Parameters))
	if err == nil {
		t.Fatal("no error")
	}
}
