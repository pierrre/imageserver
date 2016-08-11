package groupcache

import (
	"strconv"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerSize(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		benchmarkServer(b, tc.name, tc.im, 1)
	}
}

func BenchmarkServerParallelism(b *testing.B) {
	for _, p := range []int{
		1, 2, 4, 8, 16, 32, 64, 128,
	} {
		benchmarkServer(b, strconv.Itoa(p), testdata.Medium, p)
	}
}

func benchmarkServer(b *testing.B, name string, im *imageserver.Image, parallelism int) {
	b.Run(name, func(b *testing.B) {
		srv := newTestServer(
			imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return im, nil
			}),
			imageserver_cache.KeyGeneratorFunc(func(params imageserver.Params) string {
				return "test"
			}),
		)
		params := imageserver.Params{}
		b.SetParallelism(parallelism)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := srv.Get(params)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.SetBytes(int64(len(im.Data)))
	})
}
