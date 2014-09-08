package provider_test

import (
	"errors"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/provider"
	"github.com/pierrre/imageserver/testdata"
)

func TestServerInterface(t *testing.T) {
	var _ imageserver.Server = &Server{}
}

func TestServer(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": testdata.MediumFileName,
	}
	s := createTestServer()
	image, err := s.Get(parameters)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("no image")
	}
}

func TestServerErrorMissingSource(t *testing.T) {
	parameters := imageserver.Parameters{}
	s := createTestServer()
	_, err := s.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorProviderSource(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": "foobar",
	}
	s := createTestServer()
	_, err := s.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorProvider(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": "test",
	}
	s := &Server{
		Provider: ProviderFunc(func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, errors.New("error")
		}),
	}
	_, err := s.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func createTestServer() *Server {
	return &Server{
		Provider: testdata.Provider,
	}
}
