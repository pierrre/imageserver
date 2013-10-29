// Memcache cache
package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
)

// Uses Brad Fitzpatrick's Memcache client https://github.com/bradfitz/gomemcache
type MemcacheCache struct {
	Memcache *memcache_impl.Client
}

func (cache *MemcacheCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	item, err := cache.Memcache.Get(key)
	if err != nil {
		return nil, imageserver.NewCacheMissError(key, cache, err)
	}

	image, err := imageserver.NewImageUnmarshal(item.Value)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (cache *MemcacheCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	data, err := image.Marshal()
	if err != nil {
		return err
	}

	item := &memcache_impl.Item{
		Key:   key,
		Value: data,
	}
	err = cache.Memcache.Set(item)

	return err
}
