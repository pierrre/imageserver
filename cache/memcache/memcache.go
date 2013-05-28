package memcache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageproxy"
	"io"
)

type MemcacheCache struct {
	Prefix   string
	Memcache *memcache_impl.Client
}

func (cache *MemcacheCache) Get(key string) (image *imageproxy.Image, err error) {
	hashedKey := cache.hashKey(key)
	item, err := cache.Memcache.Get(hashedKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	image = &imageproxy.Image{}
	err = image.Unserialize(item.Value)
	if err != nil {
		image = nil
	}
	return
}

func (cache *MemcacheCache) Set(key string, image *imageproxy.Image) (err error) {
	serialized, err := image.Serialize()
	if err != nil {
		return
	}
	hashedKey := cache.hashKey(key)
	item := &memcache_impl.Item{
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
