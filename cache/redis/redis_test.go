package redis

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver_cache.Cache = &Cache{}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)
	defer func() {
		_ = cache.Client.Close()
	}()
	for _, expire := range []time.Duration{0, 1 * time.Minute} {
		cache.Expire = expire
		cachetest.TestGetSet(t, cache)
	}
}

func TestGetMiss(t *testing.T) {
	cache := newTestCache(t)
	defer func() {
		_ = cache.Client.Close()
	}()
	cachetest.TestGetMiss(t, cache)
}

func TestGetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress()
	defer func() {
		_ = cache.Client.Close()
	}()
	_, err := cache.Get(context.Background(), cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorAddress(t *testing.T) {
	cache := newTestCacheInvalidAddress()
	defer func() {
		_ = cache.Client.Close()
	}()
	err := cache.Set(context.Background(), cachetest.KeyValid, testdata.Medium, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)
	defer func() {
		_ = cache.Client.Close()
	}()
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	data = data[:len(data)-1]
	err = cache.setData(cachetest.KeyValid, data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cache.Get(context.Background(), cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestSetErrorMarshal(t *testing.T) {
	cache := newTestCache(t)
	defer func() {
		_ = cache.Client.Close()
	}()
	im := &imageserver.Image{
		Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
	}
	err := cache.Set(context.Background(), cachetest.KeyValid, im, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func newTestCache(tb testing.TB) *Cache {
	cache := &Cache{
		Client: newTestRedisClient("localhost:6379"),
	}
	checkTestCacheAvailable(tb, cache)
	return cache
}

func newTestCacheInvalidAddress() *Cache {
	return &Cache{
		Client: newTestRedisClient("localhost:16379"),
	}
}

func newTestRedisClient(address string) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{address},
	})
}

func checkTestCacheAvailable(tb testing.TB, cache *Cache) {
	err := cache.Client.Ping().Err()
	if err != nil {
		tb.Skip(err)
	}
}
