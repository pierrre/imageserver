package imageproxy

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	prefix   string
	memcache *memcache.Client
}

func NewMemcacheCache(prefix string, memcache *memcache.Client) *MemcacheCache {
	return &MemcacheCache{
		prefix:   prefix,
		memcache: memcache,
	}
}

func (c *MemcacheCache) Get(key string) *Image {
	return nil
}

func (c *MemcacheCache) Set(key string, image *Image) {
}
