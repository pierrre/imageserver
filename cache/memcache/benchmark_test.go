package memcache

import (
	"testing"

	"github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkGetSmall(b *testing.B) {
	benchmarkGet(b, testdata.Small)
}

func BenchmarkGetMedium(b *testing.B) {
	benchmarkGet(b, testdata.Medium)
}

func BenchmarkGetLarge(b *testing.B) {
	benchmarkGet(b, testdata.Large)
}

func benchmarkGet(b *testing.B, image *imageserver.Image) {
	cache := newTestCache(b)
	cachetest.BenchmarkGet(b, cache, 1, image) // memcached is unstable with more parallelism
}
