// Package memcache provides a Memcache Image Cache
package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// MemcacheCache represents a Memcache Image Cache
//
// It uses Brad Fitzpatrick's Memcache client https://github.com/bradfitz/gomemcache
type MemcacheCache struct {
	Client *memcache_impl.Client
}

// Get gets an Image from Memcache
func (cache *MemcacheCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	data, err := cache.getData(key)
	if err != nil {
		return nil, err
	}

	image, err := imageserver.NewImageUnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (cache *MemcacheCache) getData(key string) ([]byte, error) {
	item, err := cache.Client.Get(key)
	if err != nil {
		return nil, &imageserver_cache.MissError{Key: key}
	}

	return item.Value, nil
}

// Set sets an Image to Memcache
func (cache *MemcacheCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	data, _ := image.MarshalBinary()

	err := cache.setData(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (cache *MemcacheCache) setData(key string, data []byte) error {
	err := cache.Client.Set(&memcache_impl.Item{
		Key:   key,
		Value: data,
	})
	if err != nil {
		return err
	}

	return nil
}
