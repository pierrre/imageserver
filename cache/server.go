package cache

import (
	"encoding/hex"
	"hash"
	"io"
	"sync"

	"github.com/pierrre/imageserver"
)

// Server is a imageserver.Server implementation that supports a Cache.
//
// Steps:
//  - Generate the cache key.
//  - Get the Image from the Cache, and return it if found.
//  - Get the Image from the Server.
//  - Set the Image to the Cache.
//  - Return the Image.
type Server struct {
	imageserver.Server
	Cache        Cache
	KeyGenerator KeyGenerator
}

// Get implements imageserver.Server.
func (s *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	key := s.KeyGenerator.GetKey(params)
	im, err := s.Cache.Get(key, params)
	if err != nil {
		return nil, err
	}
	if im != nil {
		return im, nil
	}
	im, err = s.Server.Get(params)
	if err != nil {
		return nil, err
	}
	err = s.Cache.Set(key, im, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}

// KeyGenerator represents a Cache key generator.
type KeyGenerator interface {
	GetKey(imageserver.Params) string
}

// KeyGeneratorFunc is a KeyGenerator func.
type KeyGeneratorFunc func(imageserver.Params) string

// GetKey implements KeyGenerator.
func (f KeyGeneratorFunc) GetKey(params imageserver.Params) string {
	return f(params)
}

// NewParamsHashKeyGenerator returns a new KeyGenerator that hashes the Params.
func NewParamsHashKeyGenerator(newHashFunc func() hash.Hash) KeyGenerator {
	pool := &sync.Pool{
		New: func() interface{} {
			return newHashFunc()
		},
	}
	return KeyGeneratorFunc(func(params imageserver.Params) string {
		h := pool.Get().(hash.Hash)
		io.WriteString(h, params.String())
		data := h.Sum(nil)
		h.Reset()
		pool.Put(h)
		return hex.EncodeToString(data)
	})
}

// PrefixKeyGenerator is a KeyGenerator implementation that adds a prefix to the key.
type PrefixKeyGenerator struct {
	KeyGenerator
	Prefix string
}

// GetKey implements KeyGenerator.
func (g *PrefixKeyGenerator) GetKey(params imageserver.Params) string {
	return g.Prefix + g.KeyGenerator.GetKey(params)
}
