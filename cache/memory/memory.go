package memory

import (
	"fmt"
	"github.com/pierrre/imageserver"
	lru_impl "github.com/pierrre/imageserver/cache/memory/lru"
)

type MemoryCache struct {
	lru *lru_impl.LRUCache
}

func New(capacity uint64) *MemoryCache {
	return &MemoryCache{
		lru: lru_impl.NewLRUCache(capacity),
	}
}

func (cache *MemoryCache) Get(key string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	value, ok := cache.lru.Get(key)
	if !ok {
		err = fmt.Errorf("Not found")
		return
	}
	item := value.(*item)
	image = item.image
	return
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
