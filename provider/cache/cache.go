// Package cache provides a cached Image Provider
package cache

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_provider "github.com/pierrre/imageserver/provider"
)

// CacheProvider represents a cached Image Provider
type CacheProvider struct {
	Provider          imageserver_provider.Provider
	Cache             imageserver_cache.Cache
	CacheKeyGenerator CacheKeyGenerator
}

// Get returns an Image for a source
//
// It caches the image.
// The cache key used is a sha256 of the source's string representation.
func (provider *CacheProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	cacheKey := provider.CacheKeyGenerator.GetKey(source, parameters)

	image, err := provider.Cache.Get(cacheKey, parameters)
	if err == nil {
		return image, nil
	}

	image, err = provider.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	err = provider.Cache.Set(cacheKey, image, parameters)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// CacheKeyGenerator generates a Cache Key
type CacheKeyGenerator interface {
	GetKey(source interface{}, parameters imageserver.Parameters) string
}

// CacheKeyGeneratorFunc is a CacheKeyGenerator func
type CacheKeyGeneratorFunc func(source interface{}, parameters imageserver.Parameters) string

// GetKey calls the func
func (f CacheKeyGeneratorFunc) GetKey(source interface{}, parameters imageserver.Parameters) string {
	return f(source, parameters)
}

// NewSourceHashCacheKeyGeneratorFunc returns a CacheKeyGeneratorFunc that hashes the source
func NewSourceHashCacheKeyGeneratorFunc(newHashFunc func() hash.Hash) CacheKeyGeneratorFunc {
	return CacheKeyGeneratorFunc(func(source interface{}, parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, fmt.Sprint(source))
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	})
}
