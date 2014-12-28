package http

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewParamsHashETagFunc(b *testing.B) {
	params := imageserver.Params{"foo": "bar"}
	f := NewParamsHashETagFunc(sha256.New)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(params)
		}
	})
}
