// Package async provides an asynchronous cache
package async

import (
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// Cache represent an asynchronous cache
type Cache struct {
	imageserver_cache.Cache
	ErrFunc func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters)
}

// Set sets an Image to the underlying Cache using another goroutine
func (cache *Cache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	go func() {
		err := cache.Cache.Set(key, image, parameters)
		if err != nil && cache.ErrFunc != nil {
			cache.ErrFunc(err, key, image, parameters)
		}
	}()

	return nil
}
