// In memory cache
package memory

import (
	"github.com/pierrre/imageserver"
	lru_impl "github.com/pierrre/imageserver/cache/memory/lru"
)

// Uses an LRU implentation from https://github.com/youtube/vitess/blob/master/go/cache/lru_cache.go
type MemoryCache struct {
	lru *lru_impl.LRUCache
}

// The capacity is the maximum cache size (in bytes)
func New(capacity int64) *MemoryCache {
	return &MemoryCache{
		lru: lru_impl.NewLRUCache(capacity),
	}
}

func (cache *MemoryCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	value, ok := cache.lru.Get(key)
	if !ok {
		return nil, imageserver.NewCacheMissError(key, cache, nil)
	}
	item := value.(*item)
	image := item.image
	return image, nil
}

func (cache *MemoryCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	item := &item{
		image: image,
	}
	cache.lru.Set(key, item)
	return nil
}

type item struct {
	image *imageserver.Image
}

func (item *item) Size() int {
	return len(item.image.Data)
}
