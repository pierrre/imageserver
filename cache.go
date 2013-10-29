package imageserver

import (
	"fmt"
)

// Image cache interface
//
// The Get() method MUST return a CacheMissError if there is no image for the key or the cache is not available.
// The other error types will be treated as "fatal" errors.
//
// The "parameters" argument can be used for custom behavior (no-cache, expiration, ...)
type Cache interface {
	Get(key string, parameters Parameters) (*Image, error)
	Set(key string, image *Image, parameters Parameters) error
}

type CacheMissError struct {
	Key   string
	Cache Cache
	Err   error
}

func NewCacheMissError(key string, cache Cache, err error) *CacheMissError {
	return &CacheMissError{
		Key:   key,
		Cache: cache,
		Err:   err,
	}
}

func (err *CacheMissError) Error() string {
	return fmt.Sprintf("Cache miss [%s] (%s)", err.Key, err.Cache)
}
