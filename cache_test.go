package imageserver_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	. "github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
)

func TestNewCacheMissError(t *testing.T) {
	key := "foobar"
	cache := cachetest.NewCacheMap()
	previousErr := fmt.Errorf("not found")

	err := NewCacheMissError(key, cache, previousErr)
	err.Error()
}

func TestNewParametersHashCacheKeyFunc(t *testing.T) {
	f := NewParametersHashCacheKeyFunc(sha256.New)
	parameters := Parameters{
		"foo": "bar",
	}
	f(parameters)
}
