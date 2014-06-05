package imageserver_test

import (
	"crypto/sha256"
	"testing"

	. "github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestServerGet(t *testing.T) {
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

func TestServerGetWithCache(t *testing.T) {
	imageServer := createImageServer()
	imageServer.Cache = cachetest.NewMapCache()
	imageServer.CacheKeyFunc = NewParametersHashCacheKeyFunc(sha256.New)

	image, err := imageServer.Get(Parameters{
		"source": testdata.MediumFileName,
	})
	if err != nil {
		t.Fatal(err)
	}

	sameImage, err := imageServer.Get(Parameters{
		"source": testdata.MediumFileName,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !ImageEqual(image, sameImage) {
		t.Fatal("not equals")
	}
}

func TestServerGetErrorMissingSource(t *testing.T) {
	_, err := createImageServer().Get(Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerGetErrorProvider(t *testing.T) {
	_, err := createImageServer().Get(Parameters{
		"source": "foobar",
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerGetErrorProcessor(t *testing.T) {
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
