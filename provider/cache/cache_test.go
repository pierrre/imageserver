package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_provider "github.com/pierrre/imageserver/provider"
)

func TestProviderInterface(t *testing.T) {
	var _ imageserver_provider.Provider = &Provider{}
}

func TestKeyGeneratorFuncInterface(t *testing.T) {
	var _ KeyGenerator = KeyGeneratorFunc(nil)
}

func TestNewSourceHashKeyGenerator(t *testing.T) {
	g := NewSourceHashKeyGenerator(sha256.New)
	g.GetKey("foobar", imageserver.Parameters{})
}
