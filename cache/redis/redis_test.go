package redis

import (
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestGetSet(t *testing.T) {
	cache := newTestCache()
	defer cache.Close()
	checkTestRedigoAvailable(t, cache)

	for _, expire := range []time.Duration{0, 1 * time.Minute} {
		cache.Expire = expire
		cachetest.CacheTestGetSetAllImages(t, cache)
	}
}

func TestGetErrorMiss(t *testing.T) {
	cache := newTestCache()
	defer cache.Close()
	checkTestRedigoAvailable(t, cache)

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func TestGetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress()
	defer cache.Close()

	_, err := cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress()
	defer cache.Close()

	err := cache.Set(cachetest.KeyValid, testdata.Medium, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache()
	defer cache.Close()
	checkTestRedigoAvailable(t, cache)

	data, _ := testdata.Medium.MarshalBinary()
	data = data[:len(data)-1]

	err := cache.setData(cachetest.KeyValid, data)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func newTestCache() *RedisCache {
	return newTestCacheWithRedigoPool(newTestRedigoPool("localhost:6379"))
}

func newTestCacheInvalidAddress() *RedisCache {
	return newTestCacheWithRedigoPool(newTestRedigoPool("localhost:16379"))
}

func newTestCacheWithRedigoPool(pool *redigo.Pool) *RedisCache {
	return &RedisCache{
		Pool: pool,
	}
}

func newTestRedigoPool(address string) *redigo.Pool {
	return &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", address)
		},
		MaxIdle: 50,
	}
}

func checkTestRedigoAvailable(t *testing.T, cache *RedisCache) {
	conn, err := cache.Pool.Dial()
	if err != nil {
		t.Skip(err)
	}
	conn.Close()
}
