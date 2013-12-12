// Package memcache provides a Memcache Image Cache
package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
)

// MemcacheCache represents a Memcache Image Cache
//
// It uses Brad Fitzpatrick's Memcache client https://github.com/bradfitz/gomemcache
type MemcacheCache struct {
	Memcache *memcache_impl.Client
}

// Get gets an Image from Memcache
func (cache *MemcacheCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	item, err := cache.Memcache.Get(key)
	if err != nil {
		return nil, imageserver.NewCacheMissError(key, cache, err)
	}

	image, err := imageserver.NewImageUnmarshalBinary(item.Value)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// Set sets an Image to Memcache
func (cache *MemcacheCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	data, err := image.MarshalBinary()
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
