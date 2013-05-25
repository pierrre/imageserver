package imageproxy

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	Prefix   string
	Memcache *memcache.Client
}

func (cache *MemcacheCache) Get(key string) (image *Image, err error) {
	item, err := cache.Memcache.Get("a")
	if err != nil {
		fmt.Println(err)
		return
	}
	image = &Image{}
	err = image.unserialize(item.Value)
	if err != nil {
		image = nil
	}
	return
}

func (cache *MemcacheCache) Set(key string, image *Image) (err error) {
	serialized, err := image.serialize()
	if err != nil {
		return
	}
	item := &memcache.Item{
		Key:   "a",
		Value: serialized,
	}
	err = cache.Memcache.Set(item)
	return
}
