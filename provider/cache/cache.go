package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pierrre/imageserver"
	"io"
)

type CacheProvider struct {
	Cache    imageserver.Cache
	Provider imageserver.Provider
}

func (provider *CacheProvider) Get(source interface{}, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
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

func (provider *CacheProvider) getCacheKey(source interface{}) string {
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprint(source))
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
