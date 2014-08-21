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

func TestKeyGeneratorFuncInterface(t *testing.T) {
	var _ KeyGenerator = KeyGeneratorFunc(nil)
}
