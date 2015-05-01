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
func (s *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	key := s.KeyGenerator.GetKey(params)

	image, err := s.Cache.Get(key, params)
	if err == nil {
		return image, nil
	}

	image, err = s.Server.Get(params)
	if err != nil {
		return nil, err
	}

	err = s.Cache.Set(key, image, params)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// KeyGenerator generates a Cache key
type KeyGenerator interface {
	GetKey(imageserver.Params) string
}

// KeyGeneratorFunc is a KeyGenerator func
type KeyGeneratorFunc func(imageserver.Params) string

// GetKey calls the func
func (f KeyGeneratorFunc) GetKey(params imageserver.Params) string {
	return f(params)
}

// NewParamsHashKeyGenerator returns a KeyGenerator that hashes the Params
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

// PrefixKeyGenerator is a KeyGenerator that adds a prefix to the key.
type PrefixKeyGenerator struct {
	KeyGenerator
	Prefix string
}

// GetKey implements KeyGenerator.
func (g *PrefixKeyGenerator) GetKey(params imageserver.Params) string {
	return g.Prefix + g.KeyGenerator.GetKey(params)
}
