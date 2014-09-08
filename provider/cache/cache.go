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

// Provider represents an Image Provider with Cache
type Provider struct {
	imageserver_provider.Provider
	Cache        imageserver_cache.Cache
	KeyGenerator KeyGenerator
}

// Get returns an Image with Cache
func (provider *Provider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	cacheKey := provider.KeyGenerator.GetKey(source, parameters)

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

// KeyGenerator generates a Cache Key
type KeyGenerator interface {
	GetKey(source interface{}, parameters imageserver.Parameters) string
}

// KeyGeneratorFunc is a KeyGenerator func
type KeyGeneratorFunc func(source interface{}, parameters imageserver.Parameters) string

// GetKey calls the func
func (f KeyGeneratorFunc) GetKey(source interface{}, parameters imageserver.Parameters) string {
	return f(source, parameters)
}

// NewSourceHashKeyGenerator returns a KeyGenerator that hashes the source
func NewSourceHashKeyGenerator(newHashFunc func() hash.Hash) KeyGenerator {
	return KeyGeneratorFunc(func(source interface{}, parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, fmt.Sprint(source))
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	})
}
