package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewSourceHashKeyGenerator(b *testing.B) {
	parameters := imageserver.Parameters{"foo": "bar"}
	g := NewParametersHashKeyGenerator(sha256.New)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.GetKey(parameters)
		}
	})
}
