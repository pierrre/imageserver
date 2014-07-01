package cache

import (
	"testing"

	imageserver_provider "github.com/pierrre/imageserver/provider"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestCacheProviderInterface(t *testing.T) {
	var _ imageserver_provider.Provider = &CacheProvider{}
}

func TestCacheKeyGeneratorFuncInterface(t *testing.T) {
	var _ CacheKeyGenerator = CacheKeyGeneratorFunc(nil)
}
