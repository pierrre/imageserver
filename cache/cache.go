// Package cache provides an Image cache.
package cache

import "github.com/pierrre/imageserver"

// Cache represents an Image cache.
//
// The Params can be used for custom behavior (no-cache, expiration, ...).
type Cache interface {
	// Must return nil and no error if the image is not found.
	Get(key string, params imageserver.Params) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, params imageserver.Params) error
}

// IgnoreError is a Cache that ignores error from the underlying Cache.
type IgnoreError struct {
	Cache
}

// Get implements Cache.
func (c *IgnoreError) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	im, err := c.Cache.Get(key, params)
	if err != nil {
		return nil, nil
	}
	return im, nil
}

// Set implements Cache.
func (c *IgnoreError) Set(key string, image *imageserver.Image, params imageserver.Params) error {
	c.Cache.Set(key, image, params)
	return nil
}

// Async is an asynchronous Cache.
//
// The Images are set from a new goroutine.
type Async struct {
	Cache
}

// Set implements Cache.
func (a *Async) Set(key string, image *imageserver.Image, params imageserver.Params) error {
	go func() {
		a.Cache.Set(key, image, params)
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
