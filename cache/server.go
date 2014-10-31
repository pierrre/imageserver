package cache

import (
	"encoding/hex"
	"hash"
	"io"
	"sync"

	"github.com/pierrre/imageserver"
)

// Server represents a Server with Cache
type Server struct {
	imageserver.Server
	Cache        Cache
	KeyGenerator KeyGenerator
}

// Get wraps the call to the underlying Server and Get from/Set to the Cache
func (s *Server) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	key := s.KeyGenerator.GetKey(parameters)

	image, err := s.Cache.Get(key, parameters)
	if err == nil {
		return image, nil
	}

	image, err = s.Server.Get(parameters)
	if err != nil {
		return nil, err
	}

	err = s.Cache.Set(key, image, parameters)
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

// NewParametersHashKeyGenerator returns a KeyGenerator that hashes the Parameters
func NewParametersHashKeyGenerator(newHashFunc func() hash.Hash) KeyGenerator {
	pool := &sync.Pool{
		New: func() interface{} {
			return newHashFunc()
		},
	}
	return KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
		h := pool.Get().(hash.Hash)
		io.WriteString(h, parameters.String())
		data := h.Sum(nil)
		h.Reset()
		pool.Put(h)
		return hex.EncodeToString(data)
	})
}

// PrefixKeyGenerator is a KeyGenerator that adds a prefix to the key.
type PrefixKeyGenerator struct {
	KeyGenerator
	Prefix string
}

// GetKey returns the prefixed key.
func (pkg *PrefixKeyGenerator) GetKey(parameters imageserver.Parameters) string {
	return pkg.Prefix + pkg.KeyGenerator.GetKey(parameters)
}
