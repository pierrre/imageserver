package cache_test

import (
	"fmt"
	"testing"

	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
)

func TestNewCacheMissError(t *testing.T) {
	key := "foobar"
	cache := cachetest.NewMapCache()
	previousErr := fmt.Errorf("not found")

	err := NewCacheMissError(key, cache, previousErr)
	err.Error()
}
