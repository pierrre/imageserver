// Package memory provides an in-memory Image Cache
package memory

import (
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	"github.com/pierrre/lrucache"
)

// Cache represents an in-memory Image Cache
//
// It uses an LRU implementation from https://github.com/pierrre/lrucache (copy of https://github.com/youtube/vitess/tree/master/go/cache)
type Cache struct {
	lru *lrucache.LRUCache
}

// New creates a Cache
//
// capacity is the maximum cache size (in bytes)
func New(capacity int64) *Cache {
	return &Cache{
		lru: lrucache.NewLRUCache(capacity),
	}
}

// Get gets an image from the in-memory Cache
func (cache *Cache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	value, ok := cache.lru.Get(key)
	if !ok {
		return nil, &imageserver_cache.MissError{Key: key}
	}
	item := value.(*item)
	image := item.image
	return image, nil
}

// Set sets an Image to the in-memory Cache
func (cache *Cache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
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
