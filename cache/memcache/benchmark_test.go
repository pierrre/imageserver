package memcache

import (
	"testing"

	"github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkGet(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
	} {
		b.Run(tc.name, func(b *testing.B) {
			cch := newTestCache(b)
			cachetest.BenchmarkGet(b, cch, 1, tc.im) // memcached is unstable with more parallelism
		})
	}
}
