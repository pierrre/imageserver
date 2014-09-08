package redis

import (
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestInterface(t *testing.T) {
	var _ imageserver_cache.Cache = &Cache{}
}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	for _, expire := range []time.Duration{0, 1 * time.Minute} {
		cache.Expire = expire
		cachetest.CacheTestGetSetAllImages(t, cache)
	}
}

func TestGetErrorMiss(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func TestGetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress(t)
	defer cache.Close()

	_, err := cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress(t)
	defer cache.Close()

	err := cache.Set(cachetest.KeyValid, testdata.Medium, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)
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

func newTestCache(tb testing.TB) *Cache {
	cache := newTestCacheWithRedigoPool(newTestRedigoPool("localhost:6379"))
	checkTestCacheAvailable(tb, cache)
	return cache
}

func newTestCacheInvalidAddress(tb testing.TB) *Cache {
	return newTestCacheWithRedigoPool(newTestRedigoPool("localhost:16379"))
}

func newTestCacheWithRedigoPool(pool *redigo.Pool) *Cache {
	return &Cache{
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

func checkTestCacheAvailable(tb testing.TB, cache *Cache) {
	conn, err := cache.Pool.Dial()
	if err != nil {
		cache.Close()
		tb.Skip(err)
	}
	conn.Close()
}
