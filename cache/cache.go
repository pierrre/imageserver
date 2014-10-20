package cache

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

/*
Cache represents an Image cache

The Parameters can be used for custom behavior (no-cache, expiration, ...).
*/
type Cache interface {
	// Must return a MissError if it is a cache related problem.
	Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error)
	Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error
}

// MissError represents a miss error
type MissError struct {
	Key string
}

func (err *MissError) Error() string {
	return fmt.Sprintf("cache miss for key \"%s\"", err.Key)
}

// List represents a list of Image Cache
type List []Cache

/*
Get gets an Image from caches in sequential order.

If an Image is found, previous caches are filled.
*/
func (l List) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	for i, c := range l {
		image, err := c.Get(key, parameters)
		if err == nil {
			if i > 0 {
				err = l.set(key, image, parameters, i)
				if err != nil {
					return nil, err
				}
			}
			return image, nil
		}
	}

	return nil, &MissError{Key: key}
}

func (l List) set(key string, image *imageserver.Image, parameters imageserver.Parameters, indexLimit int) error {
	for i := 0; i < indexLimit; i++ {
		err := l[i].Set(key, image, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set sets the image to all caches
func (l List) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	return l.set(key, image, parameters, len(l))
}

// Async represent an asynchronous Cache
type Async struct {
	Cache
	ErrFunc func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters)
}

// Set sets an Image to the underlying Cache using another goroutine
func (a *Async) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	go func() {
		err := a.Cache.Set(key, image, parameters)
		if err != nil && a.ErrFunc != nil {
			a.ErrFunc(err, key, image, parameters)
		}
	}()

	return nil
}

// Func is an Image Cache that forwards calls to user defined functions
type Func struct {
	GetFunc func(key string, parameters imageserver.Parameters) (*imageserver.Image, error)
	SetFunc func(key string, image *imageserver.Image, parameters imageserver.Parameters) error
}

// Get forwards call to GetFunc
func (c *Func) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	return c.GetFunc(key, parameters)
}

// Set forwards call to SetFunc
func (c *Func) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	return c.SetFunc(key, image, parameters)
}
