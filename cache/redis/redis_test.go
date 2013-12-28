package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver/cache/cachetest"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	cache := createTestCache()
	defer cache.Close()

	for _, expire := range []time.Duration{0, 1 * time.Minute} {
		cache.Expire = expire
		cachetest.CacheTestGetSetAllImages(t, cache)
	}
}

func createTestCache() *RedisCache {
	return &RedisCache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", "localhost:6379")
			},
			MaxIdle: 50,
		},
	}
}
