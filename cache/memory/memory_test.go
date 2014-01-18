package memory

import (
	"github.com/pierrre/imageserver/cache/cachetest"
	"testing"
)

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
