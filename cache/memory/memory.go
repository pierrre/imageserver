// Package memory provides an in-memory imageserver/cache.Cache implementation.
package memory

import (
	"github.com/pierrre/imageserver"
	"github.com/pierrre/lrucache"
)

// Cache is an in-memory imageserver/cache.Cache implementation.
//
// It uses https://github.com/pierrre/lrucache (copy of https://github.com/youtube/vitess/tree/master/go/cache) .
type Cache struct {
	lru *lrucache.LRUCache
}

// New creates a new Cache.
//
// capacity is the maximum cache size (in bytes).
func New(capacity int64) *Cache {
	return &Cache{
		lru: lrucache.NewLRUCache(capacity),
	}
}

// Get implements imageserver/cache.Cache.
func (cache *Cache) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	value, ok := cache.lru.Get(key)
	if !ok {
		return nil, nil
	}
	item := value.(*item)
	image := item.image
	return image, nil
}

// Set implements imageserver/cache.Cache.
func (cache *Cache) Set(key string, image *imageserver.Image, params imageserver.Params) error {
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
