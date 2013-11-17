// Package prefix provides an Image Cache that adds a prefix to the key
package prefix

import (
	"github.com/pierrre/imageserver"
)

// PrefixCache represents an Image Cache that adds a prefix to the key
type PrefixCache struct {
	Prefix string
	Cache  imageserver.Cache
}

// Get adds the prefix to the key and returns the Image
func (cache *PrefixCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return cache.Cache.Get(cache.getKey(key), parameters)
}

// Set adds the prefix to the key and sets the Image
func (cache *PrefixCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	return cache.Cache.Set(cache.getKey(key), image, parameters)
}

func (cache *PrefixCache) getKey(key string) string {
	return cache.Prefix + key
}
