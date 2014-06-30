package cache

import (
	"testing"

	"github.com/pierrre/imageserver"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestCacheProviderInterface(t *testing.T) {
	var _ imageserver.Provider = &CacheProvider{}
}

func TestCacheKeyGeneratorFuncInterface(t *testing.T) {
	var _ CacheKeyGenerator = CacheKeyGeneratorFunc(nil)
}
