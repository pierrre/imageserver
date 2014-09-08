// Package list provides a list of Image Cache
package list

import (
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// Cache represents a list of Image Cache
type Cache []imageserver_cache.Cache

// Get gets an Image from caches in sequential order
//
// If an Image is found, previous caches are filled
func (cache Cache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	for i, c := range cache {
		image, err := c.Get(key, parameters)
		if err == nil {
			if i > 0 {
				err = cache.set(key, image, parameters, i)
				if err != nil {
					return nil, err
				}
			}
			return image, nil
		}
	}

	return nil, &imageserver_cache.MissError{Key: key}
}

func (cache Cache) set(key string, image *imageserver.Image, parameters imageserver.Parameters, indexLimit int) error {
	for i := 0; i < indexLimit; i++ {
		err := cache[i].Set(key, image, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set sets the image to all caches
func (cache Cache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	return cache.set(key, image, parameters, len(cache))
}
