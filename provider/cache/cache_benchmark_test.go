package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewSourceHashKeyGenerator(b *testing.B) {
	source := "foobar"
	parameters := imageserver.Parameters{}
	g := NewSourceHashKeyGenerator(sha256.New)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.GetKey(source, parameters)
		}
	})
}
