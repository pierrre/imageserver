// Package redis provides a Redis Image Cache
package redis

import (
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// RedisCache represents a Redis Image Cache
//
// It uses Gary Burd's Redis client https://github.com/garyburd/redigo
type RedisCache struct {
	Pool *redigo.Pool

	Expire time.Duration // optional
}

// Get gets an Image from Redis
func (cache *RedisCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	data, err := cache.getData(key)
	if err != nil {
		return nil, &imageserver_cache.MissError{Key: key}
	}

	image, err := imageserver.NewImageUnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (cache *RedisCache) getData(key string) ([]byte, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	return redigo.Bytes(conn.Do("GET", key))
}

// Set sets an Image to Redis
func (cache *RedisCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	data, _ := image.MarshalBinary()

	err := cache.setData(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (cache *RedisCache) setData(key string, data []byte) error {
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
func (cache *RedisCache) Close() error {
	return cache.Pool.Close()
}
