package redis

import (
	"strings"
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver_cache.Cache = &Cache{}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Pool.Close()
	for _, expire := range []time.Duration{0, 1 * time.Minute} {
		cache.Expire = expire
		cachetest.TestGetSet(t, cache)
	}
}

func TestGetMiss(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Pool.Close()
	cachetest.TestGetMiss(t, cache)
}

func TestGetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress(t)
	defer cache.Pool.Close()
	_, err := cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress(t)
	defer cache.Pool.Close()
	err := cache.Set(cachetest.KeyValid, testdata.Medium, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Pool.Close()
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	data = data[:len(data)-1]
	err = cache.setData(cachetest.KeyValid, data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestSetErrorMarshal(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Pool.Close()
	im := &imageserver.Image{
		Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
	}
	err := cache.Set(cachetest.KeyValid, im, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
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
		cache.Pool.Close()
		tb.Skip(err)
	}
	conn.Close()
}
