package redis

import (
	"testing"

	"github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkGetSizeSmall(b *testing.B) {
	benchmarkGetSize(b, testdata.Small)
}

func BenchmarkGetSizeMedium(b *testing.B) {
	benchmarkGetSize(b, testdata.Medium)
}

func BenchmarkGetSizeLarge(b *testing.B) {
	benchmarkGetSize(b, testdata.Large)
}

func BenchmarkGetSizeHuge(b *testing.B) {
	benchmarkGetSize(b, testdata.Huge)
}

func benchmarkGetSize(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, image, 1)
}

func BenchmarkGetParallelism1(b *testing.B) {
	benchmarkGetParallelism(b, 1)
}

func BenchmarkGetParallelism2(b *testing.B) {
	benchmarkGetParallelism(b, 2)
}

func BenchmarkGetParallelism4(b *testing.B) {
	benchmarkGetParallelism(b, 4)
}

func BenchmarkGetParallelism8(b *testing.B) {
	benchmarkGetParallelism(b, 8)
}

func BenchmarkGetParallelism16(b *testing.B) {
	benchmarkGetParallelism(b, 16)
}

func BenchmarkGetParallelism32(b *testing.B) {
	benchmarkGetParallelism(b, 32)
}

func BenchmarkGetParallelism64(b *testing.B) {
	benchmarkGetParallelism(b, 64)
}

func BenchmarkGetParallelism128(b *testing.B) {
	benchmarkGetParallelism(b, 128)
}

func benchmarkGetParallelism(b *testing.B, parallelism int) {
	benchmarkGet(b, testdata.Medium, parallelism)
}

func benchmarkGet(b *testing.B, image *imageserver.Image, parallelism int) {
	cache := newTestCache(b)
	defer cache.Pool.Close()
	cachetest.BenchmarkGet(b, cache, 16, image)
}
