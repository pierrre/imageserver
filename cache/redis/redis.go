package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
)

type RedisCache struct {
	Pool *redigo.Pool
}

func (redis *RedisCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	data, err := redis.getData(key, parameters)
	if err != nil {
		return
	}
	image = &imageserver.Image{}
	err = image.Unmarshal(data)
	if err != nil {
		image = nil
		return
	}
	return
}

func (redis *RedisCache) getData(key string, parameters imageserver.Parameters) (data []byte, err error) {
	conn := redis.Pool.Get()
	defer conn.Close()
	data, err = redigo.Bytes(conn.Do("GET", key))
	return
}

func (redis *RedisCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) (err error) {
	data, err := image.Marshal()
	if err != nil {
		return
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, data)
	if err != nil {
		return
	}
	return
}
