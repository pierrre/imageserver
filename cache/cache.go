package cache

import (
	"github.com/pierrre/imageserver"
)

// Cache represents an Image cache
//
// The Get() method must return a MissError if it is a cache related problem.
//
// The "parameters" argument can be used for custom behavior (no-cache, expiration, ...)
type Cache interface {
	Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error
}
