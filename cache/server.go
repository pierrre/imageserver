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
	ImageServer  imageserver.ImageServerInterface
	Cache        Cache
	KeyGenerator KeyGenerator
}

// Get wraps the call to the underlying ImageServer and Get from/Set to the Cache
func (cis *CacheImageServer) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	key := cis.KeyGenerator.GetKey(parameters)

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

// KeyGenerator generates a Cache key
type KeyGenerator interface {
	GetKey(imageserver.Parameters) string
}

// KeyGeneratorFunc is a KeyGenerator func
type KeyGeneratorFunc func(imageserver.Parameters) string

// GetKey calls the func
func (f KeyGeneratorFunc) GetKey(parameters imageserver.Parameters) string {
	return f(parameters)
}

// NewParametersHashKeyGeneratorFunc returns a KeyGeneratorFunc that hashes the Parameters
func NewParametersHashKeyGeneratorFunc(newHashFunc func() hash.Hash) KeyGeneratorFunc {
	return KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, parameters.String())
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	})
}
