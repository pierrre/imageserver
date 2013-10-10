package imageserver

import (
	"errors"
	"sync"
	"testing"
)

type size struct {
	width  int
	height int
}

type cacheMap struct {
	mutex sync.RWMutex
	data  map[string]*Image
}

func newCacheMap() *cacheMap {
	return &cacheMap{
		data: make(map[string]*Image),
	}
}

func (cache *cacheMap) Get(key string, parameters Parameters) (*Image, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	image, ok := cache.data[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return image, nil
}

func (cache *cacheMap) Set(key string, image *Image, parameters Parameters) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = image

	return nil
}

type providerSize struct{}

func (provider *providerSize) Get(source interface{}, parameters Parameters) (*Image, error) {
	size, ok := source.(size)
	if !ok {
		return nil, errors.New("source is not a size")
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
		t.Fatal("no error")
	}
}

func createServer() *Server {
	return &Server{
		Cache:     newCacheMap(),
		Provider:  new(providerSize),
		Processor: new(processorCopy),
	}
}
