package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pierrre/imageserver"
	"io"
)

type CacheProvider struct {
	Cache    imageserver.Cache
	Provider imageserver.Provider
}

func (provider *CacheProvider) Get(source string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	cacheKey := provider.getCacheKey(source)
	if image, err = provider.Cache.Get(cacheKey, parameters); err == nil {
		return
	}
	if image, err = provider.Provider.Get(source, parameters); err != nil {
		return
	}
	go func() {
		_ = provider.Cache.Set(cacheKey, image, parameters)
	}()
	return
}

func (provider *CacheProvider) getCacheKey(key string) string {
	hash := sha256.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
