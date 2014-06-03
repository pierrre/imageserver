package redis

import (
	"testing"

	"github.com/pierrre/imageserver"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

// Worker 1

func BenchmarkGetWorker1Small(b *testing.B) {
	benchmarkGetWorker1(b, testdata.Small)
}

func BenchmarkGetWorker1Medium(b *testing.B) {
	benchmarkGetWorker1(b, testdata.Medium)
}

func BenchmarkGetWorker1Large(b *testing.B) {
	benchmarkGetWorker1(b, testdata.Large)
}

func BenchmarkGetWorker1Huge(b *testing.B) {
	benchmarkGetWorker1(b, testdata.Huge)
}

func BenchmarkGetWorker1Animated(b *testing.B) {
	benchmarkGetWorker1(b, testdata.Animated)
}

func benchmarkGetWorker1(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, 1, image)
}

// Worker 2

func BenchmarkGetWorker2Small(b *testing.B) {
	benchmarkGetWorker2(b, testdata.Small)
}

func BenchmarkGetWorker2Medium(b *testing.B) {
	benchmarkGetWorker2(b, testdata.Medium)
}

func BenchmarkGetWorker2Large(b *testing.B) {
	benchmarkGetWorker2(b, testdata.Large)
}

func BenchmarkGetWorker2Huge(b *testing.B) {
	benchmarkGetWorker2(b, testdata.Huge)
}

func BenchmarkGetWorker2Animated(b *testing.B) {
	benchmarkGetWorker2(b, testdata.Animated)
}

func benchmarkGetWorker2(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, 2, image)
}

// Worker 4

func BenchmarkGetWorker4Small(b *testing.B) {
	benchmarkGetWorker4(b, testdata.Small)
}

func BenchmarkGetWorker4Medium(b *testing.B) {
	benchmarkGetWorker4(b, testdata.Medium)
}

func BenchmarkGetWorker4Large(b *testing.B) {
	benchmarkGetWorker4(b, testdata.Large)
}

func BenchmarkGetWorker4Huge(b *testing.B) {
	benchmarkGetWorker4(b, testdata.Huge)
}

func BenchmarkGetWorker4Animated(b *testing.B) {
	benchmarkGetWorker4(b, testdata.Animated)
}

func benchmarkGetWorker4(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, 4, image)
}

// Worker 8

func BenchmarkGetWorker8Small(b *testing.B) {
	benchmarkGetWorker8(b, testdata.Small)
}

func BenchmarkGetWorker8Medium(b *testing.B) {
	benchmarkGetWorker8(b, testdata.Medium)
}

func BenchmarkGetWorker8Large(b *testing.B) {
	benchmarkGetWorker8(b, testdata.Large)
}

func BenchmarkGetWorker8Huge(b *testing.B) {
	benchmarkGetWorker8(b, testdata.Huge)
}

func BenchmarkGetWorker8Animated(b *testing.B) {
	benchmarkGetWorker8(b, testdata.Animated)
}

func benchmarkGetWorker8(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, 8, image)
}

// Worker 16

func BenchmarkGetWorker16Small(b *testing.B) {
	benchmarkGetWorker16(b, testdata.Small)
}

func BenchmarkGetWorker16Medium(b *testing.B) {
	benchmarkGetWorker16(b, testdata.Medium)
}

func BenchmarkGetWorker16Large(b *testing.B) {
	benchmarkGetWorker16(b, testdata.Large)
}

func BenchmarkGetWorker16Huge(b *testing.B) {
	benchmarkGetWorker16(b, testdata.Huge)
}

func BenchmarkGetWorker16Animated(b *testing.B) {
	benchmarkGetWorker16(b, testdata.Animated)
}

func benchmarkGetWorker16(b *testing.B, image *imageserver.Image) {
	benchmarkGet(b, 16, image)
}

func benchmarkGet(b *testing.B, workerCount int, image *imageserver.Image) {
	cache := newTestCache(b)
	defer cache.Close()

	cachetest.CacheBenchmarkGet(b, cache, workerCount, image)
}
