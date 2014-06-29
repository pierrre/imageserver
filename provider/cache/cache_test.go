package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestInterfaceProvider(t *testing.T) {
	var _ imageserver.Provider = &CacheProvider{}
}

func TestCacheKeyGeneratorFuncInterface(t *testing.T) {
	var _ CacheKeyGenerator = NewSourceHashCacheKeyGeneratorFunc(sha256.New)
}
