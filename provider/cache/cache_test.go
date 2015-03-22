package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_provider "github.com/pierrre/imageserver/provider"
)

var _ imageserver_provider.Provider = &Provider{}

var _ KeyGenerator = KeyGeneratorFunc(nil)

func TestNewSourceHashKeyGenerator(t *testing.T) {
	g := NewSourceHashKeyGenerator(sha256.New)
	g.GetKey("foobar", imageserver.Params{})
}
