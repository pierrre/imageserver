// Package memcache provides a Memcache imageserver/cache.Cache implementation.
package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
)

// Cache is a Memcache imageserver/cache.Cache implementation.
//
// It uses https://github.com/bradfitz/gomemcache .
type Cache struct {
	Client *memcache_impl.Client
}

// Get implements imageserver/cache.Cache.
func (cache *Cache) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	data, err := cache.getData(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	im := new(imageserver.Image)
	err = im.UnmarshalBinaryNoCopy(data)
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (cache *Cache) getData(key string) ([]byte, error) {
	item, err := cache.Client.Get(key)
	if err != nil {
		if err == memcache_impl.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	}
	return item.Value, nil
}

// Set implements imageserver/cache.Cache.
func (cache *Cache) Set(key string, im *imageserver.Image, params imageserver.Params) error {
	data, err := im.MarshalBinary()
	if err != nil {
		return err
	}
	return cache.setData(key, data)
}

func (cache *Cache) setData(key string, data []byte) error {
	return cache.Client.Set(&memcache_impl.Item{
		Key:   key,
		Value: data,
	})
}
