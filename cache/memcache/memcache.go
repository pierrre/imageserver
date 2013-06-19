package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
)

type MemcacheCache struct {
	Memcache *memcache_impl.Client
}

func (cache *MemcacheCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	item, err := cache.Memcache.Get(key)
	if err != nil {
		return
	}
	image = &imageserver.Image{}
	if err = image.Unmarshal(item.Value); err != nil {
		image = nil
	}
	return
}

func (cache *MemcacheCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) (err error) {
	data, err := image.Marshal()
	if err != nil {
		return
	}
	item := &memcache_impl.Item{
		Key:   key,
		Value: data,
	}
	err = cache.Memcache.Set(item)
	return
}
