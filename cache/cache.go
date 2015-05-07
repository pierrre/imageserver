// Package cache provides an Image cache.
package cache

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

// Cache represents an Image cache.
//
// The Params can be used for custom behavior (no-cache, expiration, ...).
type Cache interface {
	// Must return a MissError if it is a cache related problem.
	Get(key string, params imageserver.Params) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, params imageserver.Params) error
}

// MissError is a Cache miss error.
type MissError struct {
	Key string
}

func (err *MissError) Error() string {
	return fmt.Sprintf("cache miss for key \"%s\"", err.Key)
}

// Async is an asynchronous Cache.
//
// The Images are set from a new goroutine.
//
// ErrFunc is called if there is an error.
type Async struct {
	Cache
	ErrFunc func(err error, key string, image *imageserver.Image, params imageserver.Params)
}

// Set implements Cache.
func (a *Async) Set(key string, image *imageserver.Image, params imageserver.Params) error {
	go func() {
		err := a.Cache.Set(key, image, params)
		if err != nil && a.ErrFunc != nil {
			a.ErrFunc(err, key, image, params)
		}
	}()
	return nil
}

// Func is an Image Cache that forwards calls to user defined functions
type Func struct {
	GetFunc func(key string, params imageserver.Params) (*imageserver.Image, error)
	SetFunc func(key string, image *imageserver.Image, params imageserver.Params) error
}

// Get implements Cache.
func (c *Func) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	return c.GetFunc(key, params)
}

// Set implements Cache.
func (c *Func) Set(key string, image *imageserver.Image, params imageserver.Params) error {
	return c.SetFunc(key, image, params)
}
