// Package cache provides a base for an Image cache.
package cache

import (
	"context"

	"github.com/pierrre/imageserver"
)

// Cache represents an Image cache.
//
// The Params can be used for custom behavior (no-cache, expiration, ...).
type Cache interface {
	// Get returns the Image associated to the key, or nil if not found.
	Get(ctx context.Context, key string, params imageserver.Params) (*imageserver.Image, error)

	// Set adds the Image and associate it to the key.
	Set(ctx context.Context, key string, image *imageserver.Image, params imageserver.Params) error
}

// IgnoreError is a Cache implementation that ignores error from the underlying Cache.
type IgnoreError struct {
	Cache
}

// Get implements Cache.
func (c *IgnoreError) Get(ctx context.Context, key string, params imageserver.Params) (*imageserver.Image, error) {
	im, err := c.Cache.Get(ctx, key, params)
	if err != nil {
		return nil, nil
	}
	return im, nil
}

// Set implements Cache.
func (c *IgnoreError) Set(ctx context.Context, key string, image *imageserver.Image, params imageserver.Params) error {
	_ = c.Cache.Set(ctx, key, image, params)
	return nil
}

// Async is an asynchronous Cache implementation.
//
// The Images are set from a new goroutine.
type Async struct {
	Cache
}

// Set implements Cache.
func (a *Async) Set(ctx context.Context, key string, image *imageserver.Image, params imageserver.Params) error {
	go func() {
		_ = a.Cache.Set(ctx, key, image, params)
	}()
	return nil
}

// Func is a Cache implementation that forwards calls to user defined functions
type Func struct {
	GetFunc func(ctx context.Context, key string, params imageserver.Params) (*imageserver.Image, error)
	SetFunc func(ctx context.Context, key string, image *imageserver.Image, params imageserver.Params) error
}

// Get implements Cache.
func (c *Func) Get(ctx context.Context, key string, params imageserver.Params) (*imageserver.Image, error) {
	return c.GetFunc(ctx, key, params)
}

// Set implements Cache.
func (c *Func) Set(ctx context.Context, key string, image *imageserver.Image, params imageserver.Params) error {
	return c.SetFunc(ctx, key, image, params)
}
