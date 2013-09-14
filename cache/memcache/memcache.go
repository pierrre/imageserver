package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
)

type MemcacheCache struct {
	Memcache *memcache_impl.Client
}

func (cache *MemcacheCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	item, err := cache.Memcache.Get(key)
	if err != nil {
		return nil, err
	}
	image := &imageserver.Image{}
	if err = image.Unmarshal(item.Value); err != nil {
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
