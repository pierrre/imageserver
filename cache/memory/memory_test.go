package memory

import (
	"testing"

	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
)

var _ imageserver_cache.Cache = &Cache{}

func TestGetSet(t *testing.T) {
	cache := newTestCache()
	cachetest.TestGetSet(t, cache)
}

func TestGetMiss(t *testing.T) {
	cache := newTestCache()
	cachetest.TestGetMiss(t, cache)
}

func newTestCache() *Cache {
	return New(20 * 1024 * 1024)
}
