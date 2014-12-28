// Package cache provides a cached Image Provider
package cache

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"sync"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_provider "github.com/pierrre/imageserver/provider"
)

// Provider represents an Image Provider with Cache
type Provider struct {
	imageserver_provider.Provider
	Cache        imageserver_cache.Cache
	KeyGenerator KeyGenerator
}

// Get returns an Image with Cache
func (provider *Provider) Get(source interface{}, params imageserver.Params) (*imageserver.Image, error) {
	cacheKey := provider.KeyGenerator.GetKey(source, params)

	image, err := provider.Cache.Get(cacheKey, params)
	if err == nil {
		return image, nil
	}

	image, err = provider.Provider.Get(source, params)
	if err != nil {
		return nil, err
	}

	err = provider.Cache.Set(cacheKey, image, params)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// KeyGenerator generates a Cache Key
type KeyGenerator interface {
	GetKey(source interface{}, params imageserver.Params) string
}

// KeyGeneratorFunc is a KeyGenerator func
type KeyGeneratorFunc func(source interface{}, params imageserver.Params) string

// GetKey calls the func
func (f KeyGeneratorFunc) GetKey(source interface{}, params imageserver.Params) string {
	return f(source, params)
}

// NewSourceHashKeyGenerator returns a KeyGenerator that hashes the source
func NewSourceHashKeyGenerator(newHashFunc func() hash.Hash) KeyGenerator {
	pool := &sync.Pool{
		New: func() interface{} {
			return newHashFunc()
		},
	}
	return KeyGeneratorFunc(func(source interface{}, params imageserver.Params) string {
		h := pool.Get().(hash.Hash)
		io.WriteString(h, fmt.Sprint(source))
		data := h.Sum(nil)
		h.Reset()
		pool.Put(h)
		return hex.EncodeToString(data)
	})
}
