package http

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkNewParametersHashETagFunc(b *testing.B) {
	parameters := imageserver.Parameters{"foo": "bar"}
	f := NewParametersHashETagFunc(sha256.New)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(parameters)
		}
	})
}
