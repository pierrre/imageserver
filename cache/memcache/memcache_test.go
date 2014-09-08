package memcache

import (
	"testing"

	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestInterface(t *testing.T) {
	var _ imageserver_cache.Cache = &Cache{}
}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)

	// maximum object size is only 1MB
	for _, image := range []*imageserver.Image{
		testdata.Small,
		testdata.Medium,
		testdata.Large,
	} {
		cachetest.CacheTestGetSet(t, cache, image)
	}
}

func TestGetErrorMiss(t *testing.T) {
	cache := newTestCache(t)

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func TestGetErrorServer(t *testing.T) {
	cache := newTestCacheInvalidServer(t)

	_, err := cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorServer(t *testing.T) {
	cache := newTestCacheInvalidServer(t)

	err := cache.Set(cachetest.KeyValid, testdata.Medium, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)

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
	cache := newTestCacheWithClient(newTestClient("localhost:11211"))
	checkTestCacheAvailable(tb, cache)
	return cache
}

func newTestCacheInvalidServer(tb testing.TB) *Cache {
	return newTestCacheWithClient(newTestClient("localhost:11311"))
}

func newTestCacheWithClient(client *memcache_impl.Client) *Cache {
	return &Cache{
		Client: client,
	}
}

func newTestClient(server string) *memcache_impl.Client {
	return memcache_impl.New(server)
}

func checkTestCacheAvailable(tb testing.TB, cache *Cache) {
	err := cache.Client.Set(&memcache_impl.Item{
		Key:   "ping",
		Value: []byte("ping"),
	})
	if err != nil {
		tb.Skip(err)
	}
}
