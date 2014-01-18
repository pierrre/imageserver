package memory

import (
	"github.com/pierrre/imageserver/cache/cachetest"
	"testing"
)

func TestGetSet(t *testing.T) {
	cache := createTestCache()

	cachetest.CacheTestGetSetAllImages(t, cache)
}

func TestGetErrorMiss(t *testing.T) {
	cache := createTestCache()

	cachetest.CacheTestGetErrorMiss(t, cache)
}

func createTestCache() *MemoryCache {
	return New(20 * 1024 * 1024)
}
