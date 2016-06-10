// Package redis provides a Redis imageserver/cache.Cache implementation.
package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/pierrre/imageserver"
)

// Cache is a Redis imageserver/cache.Cache implementation.
//
// It uses https://github.com/go-redis/redis .
type Cache struct {
	Client redis.UniversalClient

	// Expire is an optional expiration duration.
	Expire time.Duration
}

// Get implements imageserver/cache.Cache.
func (cache *Cache) Get(ctx context.Context, key string, params imageserver.Params) (*imageserver.Image, error) {
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
	data, err := cache.Client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// Set implements imageserver/cache.Cache.
func (cache *Cache) Set(ctx context.Context, key string, im *imageserver.Image, params imageserver.Params) error {
	data, err := im.MarshalBinary()
	if err != nil {
		return err
	}
	return cache.setData(key, data)
}

func (cache *Cache) setData(key string, data []byte) error {
	return cache.Client.Set(key, data, cache.Expire).Err()
}
