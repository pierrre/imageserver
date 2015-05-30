// Package redis provides a Redis Image Cache.
package redis

import (
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
)

// Cache is a Redis Image Cache.
//
// It uses Gary Burd's Redis client https://github.com/garyburd/redigo .
type Cache struct {
	Pool   *redigo.Pool
	Expire time.Duration // optional
}

// Get implements Cache.
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
	conn := cache.Pool.Get()
	defer conn.Close()
	data, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// Set implements Cache.
func (cache *Cache) Set(key string, im *imageserver.Image, params imageserver.Params) error {
	data, _ := im.MarshalBinary()
	return cache.setData(key, data)
}

func (cache *Cache) setData(key string, data []byte) error {
	params := []interface{}{key, data}
	if cache.Expire != 0 {
		params = append(params, "EX", strconv.Itoa(int(cache.Expire.Seconds())))
	}
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", params...)
	return err
}

// Close closes the underlying Redigo pool
func (cache *Cache) Close() error {
	return cache.Pool.Close()
}
