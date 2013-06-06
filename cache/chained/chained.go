package chained

import (
	"fmt"
	"github.com/pierrre/imageproxy"
)

type ChainedCache struct {
	Caches []imageproxy.Cache
}

func (cache *ChainedCache) Get(key string) (*imageproxy.Image, error) {
	for i, c := range cache.Caches {
		image, err := c.Get(key)
		if err == nil {
			if i > 0 {
				cache.setCaches(key, image, i)
			}
			return image, nil
		}
	}
	return nil, fmt.Errorf("Image not found in chained cache")
}

func (cache *ChainedCache) setCaches(key string, image *imageproxy.Image, indexLimit int) {
	for i := 0; i < indexLimit; i++ {
		go func(i int) {
			cache.Caches[i].Set(key, image)
		}(i)
	}
}

func (cache *ChainedCache) Set(key string, image *imageproxy.Image) (err error) {
	for _, c := range cache.Caches {
		go func(c imageproxy.Cache) {
			c.Set(key, image)
		}(c)
	}
	return
}
