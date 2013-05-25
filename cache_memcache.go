package imageproxy

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	Prefix   string
	Memcache *memcache.Client
}

func (c *MemcacheCache) Get(key string) (image *Image, err error) {
	return nil, nil
}

func (c *MemcacheCache) Set(key string, image *Image) error {
	return nil
}
