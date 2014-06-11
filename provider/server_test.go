package provider_test

import (
	"errors"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/provider"
	"github.com/pierrre/imageserver/testdata"
)

func TestProviderImageServerInterface(t *testing.T) {
	var _ imageserver.ImageServer = &ProviderImageServer{}
}

func TestImageServerGet(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": testdata.MediumFileName,
	}
	pis := createTestProviderImageServer()
	image, err := pis.Get(parameters)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("no image")
	}
}

func TestImageServerGetErrorMissingSource(t *testing.T) {
	parameters := imageserver.Parameters{}
	pis := createTestProviderImageServer()
	_, err := pis.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestImageServerGetErrorProviderSource(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": "foobar",
	}
	pis := createTestProviderImageServer()
	_, err := pis.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestImageServerGetErrorProvider(t *testing.T) {
	parameters := imageserver.Parameters{
		"source": "test",
	}
	pis := &ProviderImageServer{
		Provider: ProviderFunc(func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, errors.New("error")
		}),
	}
	_, err := pis.Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func createTestProviderImageServer() *ProviderImageServer {
	return &ProviderImageServer{
		Provider: testdata.Provider,
	}
}

func TestSourceError(t *testing.T) {
	err := NewSourceError("test")
	err.Error()
}
