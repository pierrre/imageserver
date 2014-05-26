package imageserver

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// Cache represents an Image cache
//
// The Get() method must return a CacheMissError if it is a cache related problem.
//
// The "parameters" argument can be used for custom behavior (no-cache, expiration, ...)
type Cache interface {
	Get(key string, parameters Parameters) (*Image, error)
	Set(key string, image *Image, parameters Parameters) error
}

// CacheMissError represents a cache miss error (no image found or cache not available)
type CacheMissError struct {
	Key string
	Cache
	Previous error
}

// NewCacheMissError creates a new CacheMissError
func NewCacheMissError(key string, cache Cache, previous error) *CacheMissError {
	return &CacheMissError{
		Key:      key,
		Cache:    cache,
		Previous: previous,
	}
}

func (err *CacheMissError) Error() string {
	s := fmt.Sprintf("cache miss for key \"%s\"", err.Key)
	if err.Previous != nil {
		s = fmt.Sprintf("%s (%s)", s, err.Previous)
	}
	return s
}

// NewParametersHashCacheKeyFunc returns a function that hashes the parameters and returns a Cache key
func NewParametersHashCacheKeyFunc(newHashFunc func() hash.Hash) func(parameters Parameters) string {
	return func(parameters Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, parameters.String())
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	}
}
