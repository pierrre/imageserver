package cache

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

// Cache represents an Image cache
//
// The Get() method must return a CacheMissError if it is a cache related problem.
//
// The "parameters" argument can be used for custom behavior (no-cache, expiration, ...)
type Cache interface {
	Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error
}

// MissError represents a cache miss error (no image found or cache not available)
type MissError struct {
	Key      string
	Cache    Cache
	Previous error
}

// NewMissError creates a new MissError
func NewMissError(key string, cache Cache, previous error) *MissError {
	return &MissError{
		Key:      key,
		Cache:    cache,
		Previous: previous,
	}
}

func (err *MissError) Error() string {
	s := fmt.Sprintf("cache miss for key \"%s\"", err.Key)
	if err.Previous != nil {
		s = fmt.Sprintf("%s (%s)", s, err.Previous)
	}
	return s
}
