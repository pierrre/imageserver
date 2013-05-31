package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageproxy"
)

type MemcacheCache struct {
	Memcache *memcache_impl.Client
}

func (cache *MemcacheCache) Get(key string) (image *imageproxy.Image, err error) {
	hashedKey := imageproxy.HashCacheKey(key)
	item, err := cache.Memcache.Get(hashedKey)
	if err != nil {
		return
	}
	image = &imageproxy.Image{}
	err = image.Unserialize(item.Value)
	if err != nil {
		image = nil
	}
	return
}

func (cache *MemcacheCache) Set(key string, image *imageproxy.Image) (err error) {
	serialized, err := image.Serialize()
	if err != nil {
		return
	}
	hashedKey := imageproxy.HashCacheKey(key)
	item := &memcache_impl.Item{
		Key:   hashedKey,
		Value: serialized,
	}
	err = cache.Memcache.Set(item)
	return
}
