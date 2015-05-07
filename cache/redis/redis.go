// Package redis provides a Redis Image Cache.
package redis

import (
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
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
		return nil, &imageserver_cache.MissError{Key: key}
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
	return redigo.Bytes(conn.Do("GET", key))
}

// Set implements Cache.
func (cache *Cache) Set(key string, im *imageserver.Image, params imageserver.Params) error {
	data, _ := im.MarshalBinary()
	err := cache.setData(key, data)
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) setData(key string, data []byte) error {
	params := []interface{}{key, data}
	if cache.Expire != 0 {
		params = append(params, "EX", strconv.Itoa(int(cache.Expire.Seconds())))
	}
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", params...)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the underlying Redigo pool
func (cache *Cache) Close() error {
	return cache.Pool.Close()
}
