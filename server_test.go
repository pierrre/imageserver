package imageserver

import (
	"errors"
	"testing"
)

type size struct {
	width  int
	height int
}

type providerSize struct{}

func (provider *providerSize) Get(source interface{}, parameters Parameters) (*Image, error) {
	size, ok := source.(size)
	if !ok {
		return nil, errors.New("Source is not a size")
	}
	return CreateImage(size.width, size.height), nil
}

type processorCopy struct{}

func (processor *processorCopy) Process(image *Image, parameters Parameters) (*Image, error) {
	data := make([]byte, len(image.Data))
	copy(image.Data, data)
	return &Image{
			Type: image.Type,
			Data: data,
		},
		nil
}

func TestServerGet(t *testing.T) {
	_, err := createServer().Get(Parameters{
		"source": size{
			width:  500,
			height: 400,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerGetErrorMissingSource(t *testing.T) {
	parameters := make(Parameters)
	_, err := createServer().Get(parameters)
	if err == nil {
		t.Fatal("No error")
	}
}

func createServer() *Server {
	return &Server{
		Provider:  new(providerSize),
		Processor: new(processorCopy),
	}
}
