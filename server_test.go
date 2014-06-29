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

func TestImageServerGetErrorProcessor(t *testing.T) {
	imageServer := &ImageServer{
		Provider:  testdata.Provider,
		Processor: new(errorProcessor),
	}

	_, err := imageServer.Get(Parameters{
		"source": testdata.MediumFileName,
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestImageServerFunc(t *testing.T) {
	s := ImageServerFunc(func(parameters Parameters) (*Image, error) {
		return testdata.Medium, nil
	})
	s.Get(Parameters{})
}

func createImageServer() *ImageServer {
	return &ImageServer{
		Provider:  testdata.Provider,
		Processor: new(copyProcessor),
	}
}

type copyProcessor struct{}

func (processor *copyProcessor) Process(image *Image, parameters Parameters) (*Image, error) {
	data := make([]byte, len(image.Data))
	copy(image.Data, data)
	return &Image{
			Format: image.Format,
			Data:   data,
		},
		nil
}

type errorProcessor struct{}

func (processor *errorProcessor) Process(image *Image, parameters Parameters) (*Image, error) {
	return nil, NewError("error")
}
