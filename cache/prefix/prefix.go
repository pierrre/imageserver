package prefix

import (
	"github.com/pierrre/imageproxy"
)

type PrefixCache struct {
	Prefix string
	Cache  imageproxy.Cache
}

func (cache *PrefixCache) Get(key string) (image *imageproxy.Image, err error) {
	return cache.Cache.Get(cache.Prefix + key)
}

func (cache *PrefixCache) Set(key string, image *imageproxy.Image) (err error) {
	return cache.Cache.Set(cache.Prefix+key, image)
}
