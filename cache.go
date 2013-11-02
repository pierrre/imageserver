package imageserver

import (
	"fmt"
)

// Cache represents an Image cache
//
// The Get() method MUST return a CacheMissError if it is a cache related problem.
// The other error types will be treated as "fatal" errors.
//
// The "parameters" argument can be used for custom behavior (no-cache, expiration, ...)
type Cache interface {
	Get(key string, parameters Parameters) (*Image, error)
	Set(key string, image *Image, parameters Parameters) error
}

// CacheMissError represents a cache miss error (no image found or cache not available)
type CacheMissError struct {
	Key   string
	Cache Cache
	Err   error
}

// NewCacheMissError creates a new CacheMissError
func NewCacheMissError(key string, cache Cache, err error) *CacheMissError {
	return &CacheMissError{
		Key:   key,
		Cache: cache,
		Err:   err,
	}
}

func (err *CacheMissError) Error() string {
	return fmt.Sprintf("cache miss for key %s (%s)", err.Key, err.Cache)
}
