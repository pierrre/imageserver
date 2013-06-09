package prefix

import (
	"github.com/pierrre/imageserver"
)

type PrefixCache struct {
	Prefix string
	Cache  imageserver.Cache
}

func (cache *PrefixCache) Get(key string) (image *imageserver.Image, err error) {
	return cache.Cache.Get(cache.Prefix + key)
}

func (cache *PrefixCache) Set(key string, image *imageserver.Image) (err error) {
	return cache.Cache.Set(cache.Prefix+key, image)
}
