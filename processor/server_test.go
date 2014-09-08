package processor

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestServerInterface(t *testing.T) {
	var _ imageserver.Server = &Server{}
}

func TestServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Small, nil
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return image, nil
		}),
	}
	image, err := s.Get(make(imageserver.Parameters))
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("no image")
	}
}

func TestServerErrorServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return image, nil
		}),
	}
	_, err := s.Get(make(imageserver.Parameters))
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorProcessor(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Small, nil
		}),
		Processor: ProcessorFunc(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := s.Get(make(imageserver.Parameters))
	if err == nil {
		t.Fatal("no error")
	}
}
