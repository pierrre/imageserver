// Package cachetest provides utilities for cache testing
package cachetest

import (
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"testing"
)

const (
	// KeyValid is a valid cache key (with content)
	KeyValid = "test"
	// KeyMiss is an invalid cache key (without content)
	KeyMiss = "unknown"
)

var (
	// ParametersEmpty is an empty Parameters
	ParametersEmpty = make(imageserver.Parameters)
)

// CacheTestGetSet is a helper to test cache Get()/Set()
func CacheTestGetSet(t *testing.T, cache imageserver.Cache, image *imageserver.Image) {
	err := cache.Set(KeyValid, image, ParametersEmpty)
	if err != nil {
		t.Fatal(err)
	}

	newImage, err := cache.Get(KeyValid, ParametersEmpty)
	if err != nil {
		t.Fatal(err)
	}

	if !imageserver.ImageEqual(newImage, image) {
		t.Fatal("image not equals")
	}
}

// CacheTestGetSetAllImages is a helper to test cache Get()/Set() with all images from test data
func CacheTestGetSetAllImages(t *testing.T, cache imageserver.Cache) {
	for _, image := range testdata.Images {
		CacheTestGetSet(t, cache, image)
	}
}

// CacheTestGetErrorMiss is a helper to test cache Get() with a "cache miss" error
func CacheTestGetErrorMiss(t *testing.T, cache imageserver.Cache) {
	_, err := cache.Get(KeyMiss, ParametersEmpty)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.CacheMissError); !ok {
		t.Fatal("invalid error type")
	}
}
