package imageserver_test

import (
	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"testing"
)

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

func TestServerGet(t *testing.T) {
	image, err := createServer().Get(Parameters{
		"source": "medium.jpg",
	})
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("image is nil")
	}
}

func TestServerGetErrorMissingSource(t *testing.T) {
	parameters := make(Parameters)
	_, err := createServer().Get(parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func createServer() *Server {
	return &Server{
		Cache:     newCacheMap(),
		Provider:  testdata.Provider,
		Processor: new(copyProcessor),
	}
}
