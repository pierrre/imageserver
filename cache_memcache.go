package imageproxy

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
)

type MemcacheCache struct {
	Prefix   string
	Memcache *memcache.Client
}

func (cache *MemcacheCache) Get(key string) (image *Image, err error) {
	hashedKey := cache.hashKey(key)
	item, err := cache.Memcache.Get(hashedKey)
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
	hashedKey := cache.hashKey(key)
	item := &memcache.Item{
		Key:   hashedKey,
		Value: serialized,
	}
	err = cache.Memcache.Set(item)
	return
}

func (cache *MemcacheCache) hashKey(key string) string {
	hash := md5.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	hashedKey := hex.EncodeToString(data)
	return hashedKey
}
