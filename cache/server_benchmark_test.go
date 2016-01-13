package cache

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewParamsHashKeyGenerator(b *testing.B) {
	params := imageserver.Params{"foo": "bar"}
	g := NewParamsHashKeyGenerator(sha256.New)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetKey(params)
	}
}
