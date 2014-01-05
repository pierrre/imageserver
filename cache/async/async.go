// Package async provides an asynchronous cache
package async

import (
	"github.com/pierrre/imageserver"
)

// AsyncCache represent an asynchronous cache
type AsyncCache struct {
	imageserver.Cache

	ErrFunc func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters)
}

// Set sets the Image to the underlying Cache using another goroutine
func (cache *AsyncCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	go func() {
		err := cache.Cache.Set(key, image, parameters)
		if err != nil && cache.ErrFunc != nil {
			cache.ErrFunc(err, key, image, parameters)
		}
	}()

	return nil
}
