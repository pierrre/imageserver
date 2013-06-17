package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
)

type RedisCache struct {
	Pool *redigo.Pool
}

func (redis *RedisCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	conn := redis.Pool.Get()
	defer conn.Close()
	serialized, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		return
	}
	image = &imageserver.Image{}
	err = image.Unserialize(serialized)
	if err != nil {
		image = nil
		return
	}
	return
}

func (redis *RedisCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) (err error) {
	serialized, err := image.Serialize()
	if err != nil {
		return
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, serialized)
	if err != nil {
		return
	}
	return
}
