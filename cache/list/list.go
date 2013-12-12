// Package list provides a list of Image Cache
package list

import (
	"github.com/pierrre/imageserver"
)

// ListCache represents a list of Image Cache
type ListCache []imageserver.Cache

// Get gets an Image from caches in sequential order
//
// If an Image is found, previous caches are filled
func (cache ListCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	for i, c := range cache {
		image, err := c.Get(key, parameters)

		if err == nil {
			if i > 0 {
				cache.setCaches(key, image, parameters, i)
			}
			return image, nil
		}
	}

	return nil, imageserver.NewCacheMissError(key, cache, nil)
}

func (cache ListCache) setCaches(key string, image *imageserver.Image, parameters imageserver.Parameters, indexLimit int) {
	for i := 0; i < indexLimit; i++ {
		go func(i int) {
			cache[i].Set(key, image, parameters)
		}(i)
	}
}

// Set sets the image to all caches
func (cache ListCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	for _, c := range cache {
		go func(c imageserver.Cache) {
			c.Set(key, image, parameters)
		}(c)
	}
	return nil
}
