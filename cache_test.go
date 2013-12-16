package imageserver_test

import (
	"crypto/sha256"
	"fmt"
	. "github.com/pierrre/imageserver"
	"sync"
	"testing"
)

func TestNewCacheMissError(t *testing.T) {
	key := "foobar"
	cache := newCacheMap()
	previousErr := fmt.Errorf("not found")

	err := NewCacheMissError(key, cache, previousErr)
	err.Error()
}

func TestNewParametersHashCacheKeyFunc(t *testing.T) {
	f := NewParametersHashCacheKeyFunc(sha256.New)
	parameters := Parameters{
		"foo": "bar",
	}
	f(parameters)
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
		return nil, NewCacheMissError(key, cache, nil)
	}

	return image, nil
}

func (cache *cacheMap) Set(key string, image *Image, parameters Parameters) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = image

	return nil
}

type cacheFunc struct {
	GetFunc func(key string, parameters Parameters) (*Image, error)
	SetFunc func(key string, image *Image, parameters Parameters) error
}

func (cache *cacheFunc) Get(key string, parameters Parameters) (*Image, error) {
	return cache.GetFunc(key, parameters)
}

func (cache *cacheFunc) Set(key string, image *Image, parameters Parameters) error {
	return cache.SetFunc(key, image, parameters)
}
