package cache

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

/*
Cache represents an Image cache

The Params can be used for custom behavior (no-cache, expiration, ...).
*/
type Cache interface {
	// Must return a MissError if it is a cache related problem.
	Get(key string, params imageserver.Params) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, params imageserver.Params) error
}

// MissError represents a miss error
type MissError struct {
	Key string
}

func (err *MissError) Error() string {
	return fmt.Sprintf("cache miss for key \"%s\"", err.Key)
}

// Async represent an asynchronous Cache
type Async struct {
	Cache
	ErrFunc func(err error, key string, image *imageserver.Image, params imageserver.Params)
}

// Set sets an Image to the underlying Cache using another goroutine
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

// Get forwards call to GetFunc
func (c *Func) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	return c.GetFunc(key, params)
}

// Set forwards call to SetFunc
func (c *Func) Set(key string, image *imageserver.Image, params imageserver.Params) error {
	return c.SetFunc(key, image, params)
}
