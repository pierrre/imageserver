package imageproxy

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	Prefix   string
	Memcache *memcache.Client
}

func (c *MemcacheCache) Get(key string) *Image {
	return nil
}

func (c *MemcacheCache) Set(key string, image *Image) {
}
