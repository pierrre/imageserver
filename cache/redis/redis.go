package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	"strconv"
	"time"
)

type RedisCache struct {
	Pool *redigo.Pool

	Expire time.Duration
}

func (redis *RedisCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	data, err := redis.getData(key, parameters)
	if err != nil {
		return
	}
	image = &imageserver.Image{}
	if err = image.Unmarshal(data); err != nil {
		image = nil
	}
	return
}

func (redis *RedisCache) getData(key string, parameters imageserver.Parameters) ([]byte, error) {
	conn := redis.Pool.Get()
	defer conn.Close()
	return redigo.Bytes(conn.Do("GET", key))
}

func (redis *RedisCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) (err error) {
	data, err := image.Marshal()
	if err != nil {
		return
	}
	params := []interface{}{key, data}
	if redis.Expire != 0 {
		params = append(params, "EX", strconv.Itoa(int(redis.Expire.Seconds())))
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", params...)
	return
}
