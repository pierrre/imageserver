package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestImageServerInterface(t *testing.T) {
	var _ ImageServerInterface = &ImageServer{}
}

func TestImageServerGet(t *testing.T) {
	image, err := createImageServer().Get(Parameters{
		"source": testdata.MediumFileName,
	})
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("image is nil")
	}
}

func TestImageServerGetErrorMissingSource(t *testing.T) {
	_, err := createImageServer().Get(Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestImageServerGetErrorProvider(t *testing.T) {
	_, err := createImageServer().Get(Parameters{
		"source": "foobar",
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestImageServerFuncInterface(t *testing.T) {
	var _ ImageServerInterface = ImageServerFunc(nil)
}

func TestImageServerFunc(t *testing.T) {
	s := ImageServerFunc(func(parameters Parameters) (*Image, error) {
		return testdata.Medium, nil
	})
	s.Get(Parameters{})
}

func createImageServer() *ImageServer {
	return &ImageServer{
		Provider: testdata.Provider,
	}
}
