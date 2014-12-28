package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewParamsHashKeyGenerator(b *testing.B) {
	params := imageserver.Params{"foo": "bar"}
	g := NewParamsHashKeyGenerator(sha256.New)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.GetKey(params)
		}
	})
}
