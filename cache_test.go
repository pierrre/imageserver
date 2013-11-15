package imageserver

import (
	"sync"
	"testing"
)

type CacheMap struct {
	mutex sync.RWMutex
	data  map[string]*Image
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		data: make(map[string]*Image),
	}
}

func (cache *CacheMap) Get(key string, parameters Parameters) (*Image, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	image, ok := cache.data[key]
	if !ok {
		return nil, NewCacheMissError(key, cache, nil)
	}

	return image, nil
}

func (cache *CacheMap) Set(key string, image *Image, parameters Parameters) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = image

	return nil
}

type CacheFunc struct {
	GetFunc func(key string, parameters Parameters) (*Image, error)
	SetFunc func(key string, image *Image, parameters Parameters) error
}

func (cache *CacheFunc) Get(key string, parameters Parameters) (*Image, error) {
	return cache.GetFunc(key, parameters)
}

func (cache *CacheFunc) Set(key string, image *Image, parameters Parameters) error {
	return cache.SetFunc(key, image, parameters)
}

func TestCacheMissError(t *testing.T) {
	parameters := make(Parameters)
	cache := NewCacheMap()
	_, err := cache.Get("foo", parameters)
	if err == nil {
		t.Fatal("no error")
	}
	err.Error()
}
