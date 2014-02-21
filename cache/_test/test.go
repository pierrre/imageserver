// Package _test provides utilities for cache testing
package _test

import (
	"sync"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

const (
	// KeyValid is a valid cache key (with content)
	KeyValid = "test"
	// KeyMiss is an invalid cache key (without content)
	KeyMiss = "unknown"
)

var (
	// ParametersEmpty is an empty Parameters
	ParametersEmpty = make(imageserver.Parameters)
)

// CacheTestGetSet is a helper to test cache Get()/Set()
func CacheTestGetSet(t *testing.T, cache imageserver.Cache, image *imageserver.Image) {
	err := cache.Set(KeyValid, image, ParametersEmpty)
	if err != nil {
		t.Fatal(err)
	}

	newImage, err := cache.Get(KeyValid, ParametersEmpty)
	if err != nil {
		t.Fatal(err)
	}

	if !imageserver.ImageEqual(newImage, image) {
		t.Fatal("image not equals")
	}
}

// CacheTestGetSetAllImages is a helper to test cache Get()/Set() with all images from test data
func CacheTestGetSetAllImages(t *testing.T, cache imageserver.Cache) {
	for _, image := range testdata.Images {
		CacheTestGetSet(t, cache, image)
	}
}

// CacheTestGetErrorMiss is a helper to test cache Get() with a "cache miss" error
func CacheTestGetErrorMiss(t *testing.T, cache imageserver.Cache) {
	_, err := cache.Get(KeyMiss, ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.CacheMissError); !ok {
		t.Fatal("invalid error type")
	}
}

// CacheMap is a simple Cache (it wraps a map) for tests
type CacheMap struct {
	mutex sync.RWMutex
	data  map[string]*imageserver.Image
}

// NewCacheMap creates a new CacheMap
func NewCacheMap() *CacheMap {
	return &CacheMap{
		data: make(map[string]*imageserver.Image),
	}
}

// Get gets an Image from the CacheMap
func (cache *CacheMap) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	image, ok := cache.data[key]
	if !ok {
		return nil, imageserver.NewCacheMissError(key, cache, nil)
	}

	return image, nil
}

// Set sets an Image to the CacheMap
func (cache *CacheMap) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = image

	return nil
}
