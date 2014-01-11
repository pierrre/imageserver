package memcache

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/cache/cachetest"
	"github.com/pierrre/imageserver/testdata"
	"testing"
)

func TestGetSet(t *testing.T) {
	cache := createTestCache()

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
	cache := createTestCache()

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func TestGetErrorServer(t *testing.T) {
	cache := createTestCacheInvalidServer()

	_, err := cache.Get(cachetest.KeyValid, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorServer(t *testing.T) {
	cache := createTestCacheInvalidServer()

	err := cache.Set(cachetest.KeyValid, testdata.Medium, cachetest.ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := createTestCache()

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

func createTestCache() *MemcacheCache {
	return createTestCacheWithClient(createTestClient("localhost:11211"))
}

func createTestCacheInvalidServer() *MemcacheCache {
	return createTestCacheWithClient(createTestClient("localhost:11311"))
}

func createTestCacheWithClient(client *memcache_impl.Client) *MemcacheCache {
	return &MemcacheCache{
		Client: client,
	}
}

func createTestClient(server string) *memcache_impl.Client {
	return memcache_impl.New(server)
}
