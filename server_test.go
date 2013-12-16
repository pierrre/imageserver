package imageserver_test

import (
	"crypto/sha256"
	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"reflect"
	"runtime"
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

func TestServerGetWithCache(t *testing.T) {
	server := createServer()
	server.Cache = newCacheMap()
	server.CacheKeyFunc = NewParametersHashCacheKeyFunc(sha256.New)

	image, err := server.Get(Parameters{
		"source": "medium.jpg",
	})
	if err != nil {
		t.Fatal(err)
	}

	runtime.Gosched() // TRICK: we yield this goroutine in order to fill the cache (it's set in another goroutine)

	sameImage, err := server.Get(Parameters{
		"source": "medium.jpg",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(image, sameImage) {
		t.Fatal("not equals")
	}
}

func TestServerGetErrorMissingSource(t *testing.T) {
	_, err := createServer().Get(Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerGetErrorProvider(t *testing.T) {
	_, err := createServer().Get(Parameters{
		"source": "foobar",
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func createServer() *Server {
	return &Server{
		Provider:  testdata.Provider,
		Processor: new(copyProcessor),
	}
}
