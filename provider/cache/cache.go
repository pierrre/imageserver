// Package cache provides a cached Image Provider
package cache

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/pierrre/imageserver"
)

// CacheProvider represents a cached Image Provider
type CacheProvider struct {
	imageserver.Provider
	Cache        imageserver.Cache
	CacheKeyFunc func(source interface{}, parameters imageserver.Parameters) string
}

// Get returns an Image for a source
//
// It caches the image.
// The cache key used is a sha256 of the source's string representation.
func (provider *CacheProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	cacheKey := provider.CacheKeyFunc(source, parameters)

	image, err := provider.Cache.Get(cacheKey, parameters)
	if err == nil {
		return image, nil
	}

	image, err = provider.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	provider.Cache.Set(cacheKey, image, parameters)
	// TODO handle errors properly

	return image, nil
}

// NewSourceHashCacheKeyFunc returns a function that hashes the source  and returns a Cache key
func NewSourceHashCacheKeyFunc(newHashFunc func() hash.Hash) func(source interface{}, parameters imageserver.Parameters) string {
	return func(source interface{}, parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, fmt.Sprint(source))
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	}
}
