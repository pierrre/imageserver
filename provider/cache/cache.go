// Cache provider
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pierrre/imageserver"
	"io"
)

// Cache a provider
//
// The key used is a sha256 of the source's string representation.
type CacheProvider struct {
	Cache    imageserver.Cache
	Provider imageserver.Provider
}

func (provider *CacheProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	cacheKey := provider.getCacheKey(source)

	image, err := provider.Cache.Get(cacheKey, parameters)
	if err == nil {
		return image, nil
	}

	image, err = provider.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = provider.Cache.Set(cacheKey, image, parameters)
	}()

	return image, nil
}

func (provider *CacheProvider) getCacheKey(source interface{}) string {
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprint(source))
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
