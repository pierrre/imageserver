// Package async provides an asynchronous cache
package async

import (
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// AsyncCache represent an asynchronous cache
type AsyncCache struct {
	Cache imageserver_cache.Cache

	ErrFunc func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters)
}

// Get gets an Image from the underlying Cache
func (cache *AsyncCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return cache.Cache.Get(key, parameters)
}

// Set sets an Image to the underlying Cache using another goroutine
func (cache *AsyncCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	go func() {
		err := cache.Cache.Set(key, image, parameters)
		if err != nil && cache.ErrFunc != nil {
			cache.ErrFunc(err, key, image, parameters)
		}
	}()

	return nil
}
