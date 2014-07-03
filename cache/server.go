package cache

import (
	"encoding/hex"
	"hash"
	"io"

	"github.com/pierrre/imageserver"
)

// CacheImageServer represents an Image server with Cache
//
// It wraps an ImageServer.
type CacheImageServer struct {
	ImageServer       imageserver.ImageServer
	Cache             Cache
	CacheKeyGenerator CacheKeyGenerator
}

// Get wraps the call to the underlying ImageServer and Get from/Set to the Cache
func (cis *CacheImageServer) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	key := cis.CacheKeyGenerator.GetKey(parameters)

	image, err := cis.Cache.Get(key, parameters)
	if err == nil {
		return image, nil
	}

	image, err = cis.ImageServer.Get(parameters)
	if err != nil {
		return nil, err
	}

	err = cis.Cache.Set(key, image, parameters)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// CacheKeyGenerator generates a Cache key
type CacheKeyGenerator interface {
	GetKey(imageserver.Parameters) string
}

// CacheKeyGeneratorFunc is a KeyGenerator func
type CacheKeyGeneratorFunc func(imageserver.Parameters) string

// GetKey calls the func
func (f CacheKeyGeneratorFunc) GetKey(parameters imageserver.Parameters) string {
	return f(parameters)
}

// NewParametersHashCacheKeyGeneratorFunc returns a CacheKeyGeneratorFunc that hashes the Parameters
func NewParametersHashCacheKeyGeneratorFunc(newHashFunc func() hash.Hash) CacheKeyGeneratorFunc {
	return CacheKeyGeneratorFunc(func(parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, parameters.String())
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	})
}
