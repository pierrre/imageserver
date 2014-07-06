package memory

import (
	"testing"

	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
)

func TestInterface(t *testing.T) {
	var _ imageserver_cache.Cache = &MemoryCache{}
}

func TestGetSet(t *testing.T) {
	cache := newTestCache()

	cachetest.CacheTestGetSetAllImages(t, cache)
}

func TestGetErrorMiss(t *testing.T) {
	cache := newTestCache()

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func newTestCache() *MemoryCache {
	return New(20 * 1024 * 1024)
}
