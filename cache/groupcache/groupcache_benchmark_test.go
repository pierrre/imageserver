package groupcache

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerSizeSmall(b *testing.B) {
	benchmarkServerSize(b, testdata.Small)
}

func BenchmarkServerSizeMedium(b *testing.B) {
	benchmarkServerSize(b, testdata.Medium)
}

func BenchmarkServerSizeLarge(b *testing.B) {
	benchmarkServerSize(b, testdata.Large)
}

func BenchmarkServerSizeHuge(b *testing.B) {
	benchmarkServerSize(b, testdata.Huge)
}

func benchmarkServerSize(b *testing.B, im *imageserver.Image) {
	benchmarkServer(b, im, 1)
}

func BenchmarkServerParallelism1(b *testing.B) {
	benchmarkServerParallelism(b, 1)
}

func BenchmarkServerParallelism2(b *testing.B) {
	benchmarkServerParallelism(b, 2)
}

func BenchmarkServerParallelism4(b *testing.B) {
	benchmarkServerParallelism(b, 4)
}

func BenchmarkServerParallelism8(b *testing.B) {
	benchmarkServerParallelism(b, 8)
}

func BenchmarkServerParallelism16(b *testing.B) {
	benchmarkServerParallelism(b, 16)
}

func BenchmarkServerParallelism32(b *testing.B) {
	benchmarkServerParallelism(b, 32)
}

func BenchmarkServerParallelism64(b *testing.B) {
	benchmarkServerParallelism(b, 64)
}

func BenchmarkServerParallelism128(b *testing.B) {
	benchmarkServerParallelism(b, 128)
}

func benchmarkServerParallelism(b *testing.B, parallelism int) {
	benchmarkServer(b, testdata.Medium, parallelism)
}

func benchmarkServer(b *testing.B, im *imageserver.Image, parallelism int) {
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
}
