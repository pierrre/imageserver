// Prefix cache
package prefix

import (
	"github.com/pierrre/imageserver"
)

// Concatenate a prefix to the cache key
type PrefixCache struct {
	Prefix string
	Cache  imageserver.Cache
}

func (cache *PrefixCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return cache.Cache.Get(cache.getKey(key), parameters)
}

func (cache *PrefixCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	return cache.Cache.Set(cache.getKey(key), image, parameters)
}

func (cache *PrefixCache) getKey(key string) string {
	return cache.Prefix + key
}
