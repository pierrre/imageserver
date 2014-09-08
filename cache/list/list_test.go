package list

import (
	"testing"

	imageserver_cache "github.com/pierrre/imageserver/cache"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestInterface(t *testing.T) {
	var _ imageserver_cache.Cache = Cache{}
}
