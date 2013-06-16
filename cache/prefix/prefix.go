package prefix

import (
	"github.com/pierrre/imageserver"
)

type PrefixCache struct {
	Prefix string
	Cache  imageserver.Cache
}

func (cache *PrefixCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	return cache.Cache.Get(cache.Prefix+key, parameters)
}

func (cache *PrefixCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) (err error) {
	return cache.Cache.Set(cache.Prefix+key, image, parameters)
}
