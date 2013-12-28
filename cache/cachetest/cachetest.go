// Package cachetest provides utilities for cache testing
package cachetest

import (
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"reflect"
	"testing"
)

// CacheTestGetSet is a helper to test cache Get()/Set()
func CacheTestGetSet(t *testing.T, cache imageserver.Cache, image *imageserver.Image) {
	key := "test"
	parameters := make(imageserver.Parameters)

	err := cache.Set(key, image, parameters)
	if err != nil {
		t.Fatal(err)
	}

	newImage, err := cache.Get(key, parameters)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(newImage, image) {
		t.Fatal("image not equals")
	}
}

// CacheTestGetSetAllImages is a helper to test cache Get()/Set() with all images from test data
func CacheTestGetSetAllImages(t *testing.T, cache imageserver.Cache) {
	for _, image := range testdata.Images {
		CacheTestGetSet(t, cache, image)
	}
}
