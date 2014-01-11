package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver/cache/cachetest"
	"github.com/pierrre/imageserver/testdata"
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

func TestGetErrorMiss(t *testing.T) {
	cache := createTestCache()
	defer cache.Close()

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func TestGetErrorAddress(t *testing.T) {
	cache := createTestCacheInvalidAddress()
	defer cache.Close()

	_, err := cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorAddress(t *testing.T) {
	cache := createTestCacheInvalidAddress()
	defer cache.Close()

	err := cache.Set(cachetest.KeyValid, testdata.Medium, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := createTestCache()
	defer cache.Close()

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

func createTestCache() *RedisCache {
	return createTestCacheWithRedigoPool(createTestRedigoPool("localhost:6379"))
}

func createTestCacheInvalidAddress() *RedisCache {
	return createTestCacheWithRedigoPool(createTestRedigoPool("localhost:16379"))
}

func createTestCacheWithRedigoPool(pool *redigo.Pool) *RedisCache {
	return &RedisCache{
		Pool: pool,
	}
}

func createTestRedigoPool(address string) *redigo.Pool {
	return &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", address)
		},
		MaxIdle: 50,
	}
}
