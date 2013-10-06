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

func TestServerGetSource(t *testing.T) {
	server := &Server{
		Provider: new(providerSize),
	}
	_, err := server.getSource(Parameters{
		"source": size{
			width:  500,
			height: 400,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerGetSourceErrorMissingSource(t *testing.T) {
	server := &Server{}
	parameters := make(Parameters)
	_, err := server.getSource(parameters)
	if err == nil {
		t.Fatal("No error")
	}
}
