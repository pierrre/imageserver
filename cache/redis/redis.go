// Package redis provides a Redis Image Cache
package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	"strconv"
	"time"
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
	data, err := cache.getData(key, parameters)
	if err != nil {
		return nil, imageserver.NewCacheMissError(key, cache, err)
	}

	image, err := imageserver.NewImageUnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (cache *RedisCache) getData(key string, parameters imageserver.Parameters) ([]byte, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	return redigo.Bytes(conn.Do("GET", key))
}

// Set sets an Image to Redis
func (cache *RedisCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	data, err := image.MarshalBinary()
	if err != nil {
		return err
	}

	params := []interface{}{key, data}

	if cache.Expire != 0 {
		params = append(params, "EX", strconv.Itoa(int(cache.Expire.Seconds())))
	}

	conn := cache.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", params...)
	if err != nil {
		return err
	}

	return nil
}

func (cache *RedisCache) Close() error {
	return cache.Pool.Close()
}
